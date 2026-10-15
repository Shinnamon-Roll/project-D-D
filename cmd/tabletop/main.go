// Command tabletop is the TUI launcher: a welcome screen with a start button, then a live dashboard
// (server, room, logs, connections) with a reload button that rebuilds and restarts.
// Design: docs/design/tui.md. The headless server is cmd/server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"projectdnd/internal/room"
	"projectdnd/internal/server"
)

var version = "dev" // set with -ldflags "-X main.version=..."

const (
	maxLines        = 1000
	shutdownTimeout = 5 * time.Second
)

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "listen address (0.0.0.0 = reachable from LAN; 127.0.0.1 = this machine only)")
	autostart := flag.Bool("autostart", false, "start the server right away, skipping the welcome screen (set by reload)")
	flag.Parse()

	ch := make(chan logMsg, 1024)
	log := slog.New(&teaHandler{ch: ch})
	m := newModel(*addr, version, func() (runner, error) { return server.Start(*addr, "uploads", log) })
	m.log, m.autostart, m.build = log, *autostart, makeBuild
	wd, _ := os.Getwd()       // "" on error: findRoot reports it
	exe, _ := os.Executable() // "" on error: falls back to wd
	m.root, m.noReload = findRoot(exe, wd)
	p := tea.NewProgram(m)
	go func() {
		for msg := range ch {
			p.Send(msg) // no-op once the program has exited
		}
	}()
	_, err := p.Run()
	m.shutdownIfRunning() // before os.Exit: a panic, kill or external SIGINT skips the q path
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if m.reload { // build passed and the server is down: become the new binary
		err := reexec(m.root, os.Args[1:])
		fmt.Fprintln(os.Stderr, "reload:", err)
		os.Exit(1)
	}
}

// shutdownIfRunning gracefully stops a server the q path did not already stop.
func (m *model) shutdownIfRunning() {
	if m.cancel != nil {
		m.cancel() // stop a build still running
	}
	if m.srv == nil || m.stopped {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	m.srv.Shutdown(ctx)
	m.stopped = true
}

// runner is the part of *server.Server the TUI uses (faked in tests).
type runner interface {
	Addr() string
	Clients() int
	Summary() room.Summary
	Shutdown(context.Context) error
}

type screen int

const (
	welcome screen = iota
	starting
	running
)

type (
	startedMsg struct {
		srv runner
		err error
	}
	stoppedMsg   struct{}
	tickMsg      time.Time
	buildDoneMsg struct{ err error }
)

type model struct {
	screen   screen
	addr     string
	version  string
	start    func() (runner, error)
	srv      runner
	ips      []net.IP // LAN IPv4s, for the links phones can open
	err      error    // last start failure
	since    time.Time
	stopping bool
	stopped  bool
	follow   bool
	lines    []string // rendered log lines, ring buffer, newest last
	spin     spinner.Model
	vp       viewport.Model
	w, h     int

	// reload
	log       *slog.Logger // build output goes here, so it reaches the log panel
	autostart bool
	root      string // repo root; "" = reload disabled
	noReload  string // why reload is disabled
	build     buildFunc
	confirm   time.Time // reload asked while clients are connected; a second r before this confirms
	building  bool
	cancel    context.CancelFunc // cancels the running build
	reload    bool               // build passed: main re-execs after the program exits

	// dashboard samples, taken on each tick
	sum             room.Summary
	clients         int
	connHist        []int // clients per second
	evHist          []int // events per 10s
	evNow, ticks    int
	events, rejects int // counted from log records
	warns           int // warn+error lines since the last clear
	heap            uint64
	goroutines      int
}

func newModel(addr, version string, start func() (runner, error)) *model {
	m := &model{
		addr: addr, version: version, start: start, follow: true, ips: server.LANIPs(),
		spin: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(accent)),
		vp:   viewport.New(),
		log:  slog.New(slog.DiscardHandler), noReload: "not set up",
	}
	m.vp.SoftWrap = true // long errors stay readable
	return m
}

func (m *model) Init() tea.Cmd {
	if m.autostart {
		return m.startCmd()
	}
	return nil
}

// boundAddr is the server's bound address once running, else the configured one.
func (m *model) boundAddr() string {
	if m.srv != nil {
		return m.srv.Addr()
	}
	return m.addr
}

// url is the link for this machine; an all-interfaces bind (0.0.0.0, ::, "") opens via 127.0.0.1.
func (m *model) url() string {
	addr := m.boundAddr()
	if host, port, err := net.SplitHostPort(addr); err == nil {
		if ip := net.ParseIP(host); host == "" || ip != nil && ip.IsUnspecified() {
			addr = net.JoinHostPort("127.0.0.1", port)
		}
	}
	return "http://" + addr
}

// lan is the LAN link(s) other devices can open, or "" when bound to loopback only.
func (m *model) lan() string {
	return strings.Join(server.LANURLs(m.boundAddr(), m.ips), " · ")
}

func (m *model) startCmd() tea.Cmd {
	m.screen, m.err = starting, nil
	start := m.start
	return tea.Batch(m.spin.Tick, func() tea.Msg {
		srv, err := start()
		return startedMsg{srv, err}
	})
}

func (m *model) stopCmd() tea.Cmd {
	if m.srv == nil {
		return tea.Quit
	}
	if m.stopping {
		return nil
	}
	if m.cancel != nil {
		m.cancel()
	}
	m.stopping = true
	srv := m.srv
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		srv.Shutdown(ctx)
		return stoppedMsg{}
	}
}

// reloadCmd starts the build, or first asks for confirmation when clients would lose the room.
func (m *model) reloadCmd() tea.Cmd {
	switch {
	case m.building || m.stopping:
		return nil
	case m.root == "":
		m.addLog(logMsg{time.Now(), slog.LevelWarn, "reload disabled: " + m.noReload})
		return nil
	case m.srv.Clients() > 0 && !time.Now().Before(m.confirm):
		m.confirm = time.Now().Add(confirmTimeout)
		m.addLog(logMsg{time.Now(), slog.LevelWarn, "reload: ห้องจะรีเซ็ต (ทุกคนหลุด) กด r อีกครั้งภายใน 5 วิ เพื่อยืนยัน"})
		return nil
	}
	m.confirm, m.building = time.Time{}, true
	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	m.cancel = cancel
	m.addLog(logMsg{time.Now(), slog.LevelInfo, "reload: building in " + m.root})
	build, root, log := m.build, m.root, m.log
	return tea.Batch(m.spin.Tick, func() tea.Msg {
		return buildDoneMsg{build(ctx, root, func(s string) { log.Info("build │ " + ansi.Strip(s)) })}
	})
}

// sample polls the server for the dashboard; runs on each 1s tick.
func (m *model) sample() {
	m.clients, m.sum = m.srv.Clients(), m.srv.Summary()
	m.connHist = push(m.connHist, m.clients)
	if m.ticks++; m.ticks%10 == 0 {
		m.evHist, m.evNow = push(m.evHist, m.evNow), 0
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.heap, m.goroutines = ms.HeapAlloc, runtime.NumGoroutine()
}

// addLog appends a record to the log panel and counts it for the dashboard.
func (m *model) addLog(l logMsg) {
	switch {
	case strings.HasPrefix(l.text, "event rejected"):
		m.rejects++
		m.evNow++
	case strings.HasPrefix(l.text, "event "):
		m.events++
		m.evNow++
	}
	if l.level >= slog.LevelWarn {
		m.warns++
	}
	m.lines = append(m.lines, l.render())
	if len(m.lines) > maxLines {
		m.lines = m.lines[len(m.lines)-maxLines:]
	}
	m.refresh()
}

func tick() tea.Cmd { return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }) }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		w, h := m.logSize()
		m.vp.SetWidth(max(w-4, 1)) // border + 1 space each side
		m.vp.SetHeight(max(h-2, 1))
		m.refresh()
		return m, nil
	case startedMsg:
		if msg.err != nil {
			m.screen, m.err = welcome, msg.err
			return m, nil
		}
		m.screen, m.srv, m.since = running, msg.srv, time.Now()
		m.sample()
		return m, tick()
	case stoppedMsg:
		m.stopped = true
		return m, tea.Quit
	case tickMsg:
		if m.screen == running && !m.stopping {
			m.sample()
			return m, tick()
		}
		return m, nil
	case buildDoneMsg:
		m.building = false
		m.cancel()
		if m.stopping { // q during the build wins: quit, no reload
			return m, nil
		}
		if msg.err != nil {
			m.addLog(logMsg{time.Now(), slog.LevelError, "reload: build failed, server keeps running err=" + fmt.Sprintf("%q", msg.err.Error())})
			return m, nil
		}
		m.reload = true
		m.addLog(logMsg{time.Now(), slog.LevelInfo, "reload: build ok, restarting…"})
		return m, m.stopCmd()
	case logMsg:
		m.addLog(msg)
		return m, nil
	case spinner.TickMsg:
		if m.screen != starting && !m.building {
			return m, nil
		}
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	case tea.MouseClickMsg:
		if msg.Button != tea.MouseLeft {
			return m, nil
		}
		switch {
		case m.screen == welcome && hit(m.welcomeView(), m.buttonLabel(), msg.X, msg.Y, 3, 1):
			return m, m.startCmd()
		case m.screen == running && hit(m.runningView(), reloadLabel, msg.X, msg.Y, 1, 0):
			return m, m.reloadCmd()
		}
		return m, nil
	case tea.MouseWheelMsg:
		if m.screen == running {
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			m.follow = m.vp.AtBottom()
			return m, cmd
		}
		return m, nil
	case tea.KeyPressMsg:
		return m, m.key(msg)
	}
	return m, nil
}

func (m *model) key(msg tea.KeyPressMsg) tea.Cmd {
	k := msg.String()
	if k == "q" || k == "ctrl+c" {
		return m.stopCmd()
	}
	switch m.screen {
	case welcome:
		if k == "enter" || k == "space" {
			return m.startCmd()
		}
	case running:
		switch k {
		case "o":
			return openBrowser(m.url())
		case "f":
			m.follow = !m.follow
			if m.follow {
				m.vp.GotoBottom()
			}
		case "c":
			m.lines, m.warns = nil, 0
			m.refresh()
		case "r":
			return m.reloadCmd()
		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			m.follow = m.vp.AtBottom()
			return cmd
		}
	}
	return nil
}

func openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		name := "xdg-open"
		if runtime.GOOS == "darwin" {
			name = "open"
		}
		if err := exec.Command(name, url).Start(); err != nil {
			return logMsg{time.Now(), slog.LevelWarn, "cannot open browser err=" + fmt.Sprintf("%q", err.Error())}
		}
		return nil
	}
}

// refresh re-renders the log viewport content.
// shortcut: re-sets all <=1000 lines per record (~0.5ms); batch on a tick if logs get bursty.
func (m *model) refresh() {
	m.vp.SetContentLines(m.lines)
	if m.follow {
		m.vp.GotoBottom()
	}
}

// ---- view ----

var (
	mint      = lipgloss.Color("#5EEAD4")
	accent    = lipgloss.NewStyle().Foreground(mint)
	dim       = lipgloss.NewStyle().Faint(true)
	warn      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	bad       = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	good      = lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399"))
	bold      = lipgloss.NewStyle().Bold(true)
	card      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238")).Padding(1, 6)
	buttonOn  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mint).Padding(0, 2).Bold(true)
	buttonOff = buttonOn.BorderForeground(lipgloss.Color("8")).Bold(false)
)

const logo = "▀█▀ ▄▀█ █▄▄ █   █▀▀ ▀█▀ █▀█ █▀█\n" +
	" █  █▀█ █▄█ █▄▄ ██▄  █  █▄█ █▀▀"

func (m *model) buttonLabel() string {
	if m.err != nil {
		return "↻  ลองใหม่"
	}
	return "▶  เริ่มเซิร์ฟเวอร์"
}

func (m *model) welcomeView() string {
	var btn string
	if m.screen == starting {
		btn = buttonOff.Render(m.spin.View() + " กำลังเริ่ม…")
	} else {
		btn = buttonOn.Render(accent.Render(m.buttonLabel()))
	}
	parts := []string{
		dim.Render("ยินดีต้อนรับสู่"), "",
		accent.Bold(true).Render(logo),
		"โต๊ะเล่นเกมเล่าเรื่องของแก๊งเพื่อน", "",
		btn,
	}
	if m.err != nil {
		parts = append(parts, bad.Render("เริ่มไม่ได้: "+m.err.Error()))
	}
	parts = append(parts, "", dim.Render(m.versionTag()+" · "+m.url()))
	if lan := m.lan(); lan != "" {
		parts = append(parts, dim.Render("LAN "+lan))
	}
	body := lipgloss.JoinVertical(lipgloss.Center,
		card.Render(lipgloss.JoinVertical(lipgloss.Center, parts...)),
		dim.Render("enter/space เริ่ม · q ออก"))
	if m.w == 0 {
		return body
	}
	return lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, body)
}

func (m *model) versionTag() string {
	if m.version == "dev" {
		return "dev"
	}
	return "v" + m.version
}

func (m *model) View() tea.View {
	var v tea.View
	if m.screen == running {
		v = tea.NewView(m.runningView())
	} else {
		v = tea.NewView(m.welcomeView())
	}
	v.MouseMode = tea.MouseModeCellMotion // clickable buttons; most terminals still select text with shift+drag
	v.AltScreen = true
	v.WindowTitle = "tabletop"
	return v
}

func (l logMsg) render() string {
	lvl := fmt.Sprintf("%-5s", l.level.String())
	text := lvl + " " + l.text
	switch {
	case l.level < slog.LevelInfo:
		text = dim.Render(text)
	case l.level >= slog.LevelError:
		text = bad.Render(text)
	case l.level >= slog.LevelWarn:
		text = warn.Render(text)
	}
	return dim.Render(l.time.Format("15:04:05")) + " " + text
}
