// Command tabletop is the TUI launcher: a welcome screen with a start button, then live server logs.
// Design: docs/design/tui.md. The headless server is cmd/server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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

	"projectdnd/internal/server"
)

var version = "dev" // set with -ldflags "-X main.version=..."

const (
	maxLines        = 1000
	shutdownTimeout = 5 * time.Second
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()

	ch := make(chan logMsg, 1024)
	log := slog.New(&teaHandler{ch: ch})
	m := newModel(*addr, version, func() (runner, error) { return server.Start(*addr, "uploads", log) })
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
}

// shutdownIfRunning gracefully stops a server the q path did not already stop.
func (m *model) shutdownIfRunning() {
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
	stoppedMsg struct{}
	tickMsg    time.Time
)

type model struct {
	screen   screen
	addr     string
	version  string
	start    func() (runner, error)
	srv      runner
	err      error // last start failure
	since    time.Time
	stopping bool
	stopped  bool
	follow   bool
	lines    []string // rendered log lines, ring buffer, newest last
	spin     spinner.Model
	vp       viewport.Model
	w, h     int
}

func newModel(addr, version string, start func() (runner, error)) *model {
	m := &model{
		addr: addr, version: version, start: start, follow: true,
		spin: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(accent)),
		vp:   viewport.New(),
	}
	m.vp.SoftWrap = true // long errors stay readable
	return m
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) url() string {
	if m.srv != nil {
		return "http://" + m.srv.Addr()
	}
	return "http://" + m.addr
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
	m.stopping = true
	srv := m.srv
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		srv.Shutdown(ctx)
		return stoppedMsg{}
	}
}

func tick() tea.Cmd { return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }) }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.vp.SetWidth(msg.Width)
		m.vp.SetHeight(max(msg.Height-4, 1)) // header, 2 rules, footer
		m.refresh()
		return m, nil
	case startedMsg:
		if msg.err != nil {
			m.screen, m.err = welcome, msg.err
			return m, nil
		}
		m.screen, m.srv, m.since = running, msg.srv, time.Now()
		return m, tick()
	case stoppedMsg:
		m.stopped = true
		return m, tea.Quit
	case tickMsg:
		if m.screen == running && !m.stopping {
			return m, tick()
		}
		return m, nil
	case logMsg:
		m.lines = append(m.lines, msg.render())
		if len(m.lines) > maxLines {
			m.lines = m.lines[len(m.lines)-maxLines:]
		}
		m.refresh()
		return m, nil
	case spinner.TickMsg:
		if m.screen != starting {
			return m, nil
		}
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	case tea.MouseClickMsg:
		if m.screen == welcome && msg.Button == tea.MouseLeft && m.onButton(msg.X, msg.Y) {
			return m, m.startCmd()
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
			m.lines = nil
			m.refresh()
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
	accent    = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	dim       = lipgloss.NewStyle().Faint(true)
	warn      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	bad       = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	good      = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	bold      = lipgloss.NewStyle().Bold(true)
	card      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("208")).Padding(1, 6)
	buttonOn  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("208")).Padding(0, 2).Bold(true)
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
	body := lipgloss.JoinVertical(lipgloss.Center,
		card.Render(lipgloss.JoinVertical(lipgloss.Center, parts...)),
		dim.Render("enter/space เริ่ม · q ออก"))
	if m.w == 0 {
		return body
	}
	return lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, body)
}

// onButton hit-tests a click against the rendered button (its label row ±1 for the border).
func (m *model) onButton(x, y int) bool {
	label := m.buttonLabel()
	for i, line := range strings.Split(ansi.Strip(m.welcomeView()), "\n") {
		j := strings.Index(line, label)
		if j < 0 {
			continue
		}
		col := ansi.StringWidth(line[:j])
		return y >= i-1 && y <= i+1 && x >= col-3 && x < col+ansi.StringWidth(label)+3
	}
	return false
}

func (m *model) versionTag() string {
	if m.version == "dev" {
		return "dev"
	}
	return "v" + m.version
}

func (m *model) runningView() string {
	status := good.Render("● RUNNING")
	if m.stopping {
		status = warn.Render("● STOPPING")
	}
	up := time.Since(m.since).Truncate(time.Second)
	header := fmt.Sprintf(" %s  %s   uptime %02d:%02d:%02d   clients %d   %s",
		status, bold.Render(m.url()), int(up.Hours()), int(up.Minutes())%60, int(up.Seconds())%60,
		m.srv.Clients(), dim.Render(m.versionTag()))
	follow := "off"
	if m.follow {
		follow = "on"
	}
	footer := dim.Render(" o เปิดเบราว์เซอร์ · ↑↓ pgup/pgdn เลื่อน · f ตามล่าสุด " + follow + " · c ล้าง · q หยุด+ออก")
	rule := dim.Render(strings.Repeat("─", max(m.w, 1)))
	return strings.Join([]string{header, rule, m.vp.View(), rule, footer}, "\n")
}

func (m *model) View() tea.View {
	var v tea.View
	if m.screen == running {
		v = tea.NewView(m.runningView()) // no mouse mode here, so log text stays selectable
	} else {
		v = tea.NewView(m.welcomeView())
		v.MouseMode = tea.MouseModeCellMotion
	}
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
	return " " + dim.Render(l.time.Format("15:04:05")) + " " + text
}
