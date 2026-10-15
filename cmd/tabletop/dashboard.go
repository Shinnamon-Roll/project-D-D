package main

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"projectdnd/internal/room"
	"projectdnd/internal/server"
)

// Running dashboard (docs/design/tui.md "v0.6 dashboard"): SERVER + ROOM | LOGS | CONNECTIONS + ACTIVITY + SYSTEM.

const (
	side        = 34 // side column width, border included
	wideMin     = 130
	midMin      = 90
	reloadLabel = "⟳ RELOAD"
	histLen     = 60
)

var (
	edge  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	title = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

// logSize is the LOGS panel size (border included) for the current layout.
func (m *model) logSize() (w, h int) {
	switch {
	case m.w >= wideMin:
		return m.w - 2*side, m.h - 1
	case m.w >= midMin:
		return m.w - side, m.h - 1
	}
	return m.w, m.h - 2 // one-line status on top
}

// panel draws a thin rounded box of exactly w×h cells, title in the top border, note right-aligned there.
func panel(name, note string, body []string, w, h int) string {
	if w < 5 || h < 2 {
		return ""
	}
	t := ansi.Truncate(" "+name+" ", w-2, "")
	if note != "" {
		note = " " + note + " "
	}
	fill := w - 2 - ansi.StringWidth(t) - ansi.StringWidth(note)
	if fill < 1 {
		note, fill = "", w-2-ansi.StringWidth(t)
	}
	in := w - 4
	lines := []string{edge.Render("╭") + title.Render(t) + edge.Render(strings.Repeat("─", fill)) + title.Render(note) + edge.Render("╮")}
	for i := range h - 2 {
		l := ""
		if i < len(body) {
			l = ansi.Truncate(body[i], in, "…")
		}
		lines = append(lines, edge.Render("│")+" "+l+strings.Repeat(" ", max(in-ansi.StringWidth(l), 0))+" "+edge.Render("│"))
	}
	lines = append(lines, edge.Render("╰"+strings.Repeat("─", w-2)+"╯"))
	return strings.Join(lines, "\n")
}

func (m *model) runningView() string {
	lw, lh := m.logSize()
	note := fmt.Sprintf("%d · %d warn", len(m.lines), m.warns)
	if m.w < wideMin && m.w >= midMin { // right column hidden: its key numbers move here
		note += fmt.Sprintf(" · %d online · ev %d · rej %d", m.clients, m.events, m.rejects)
	}
	logs := panel("LOGS", note, strings.Split(m.vp.View(), "\n"), lw, lh)
	var body string
	switch {
	case m.w >= wideMin:
		body = lipgloss.JoinHorizontal(lipgloss.Top, m.leftColumn(lh), logs, m.rightColumn(lh))
	case m.w >= midMin:
		body = lipgloss.JoinHorizontal(lipgloss.Top, m.leftColumn(lh), logs)
	default:
		body = m.statusLine() + "\n" + logs
	}
	follow := "off"
	if m.follow {
		follow = "on"
	}
	footer := dim.Render(ansi.Truncate(" o เปิดเว็บ · ↑↓ pgup/pgdn เลื่อน · f ตามล่าสุด "+follow+" · c ล้าง · r reload · q หยุด+ออก", max(m.w, 1), "…"))
	return body + "\n" + footer
}

func (m *model) statusText() string {
	switch {
	case m.stopping:
		return warn.Render("● STOPPING")
	case m.building:
		return accent.Render(m.spin.View() + " BUILDING…")
	}
	return good.Render("● RUNNING")
}

func (m *model) uptime() string {
	up := time.Since(m.since).Truncate(time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", int(up.Hours()), int(up.Minutes())%60, int(up.Seconds())%60)
}

func (m *model) reloadButton() string {
	switch {
	case m.root == "":
		return dim.Render(reloadLabel + " off") // reason goes to the log on r
	case m.building || m.stopping:
		return dim.Render(reloadLabel)
	case time.Now().Before(m.confirm):
		return warn.Bold(true).Render(reloadLabel + "? ห้องจะรีเซ็ต กด r อีกครั้ง")
	}
	return accent.Bold(true).Render(reloadLabel) + dim.Render(" (r)")
}

func (m *model) leftColumn(h int) string {
	srv := []string{m.statusText() + "  " + dim.Render(m.versionTag()), bold.Render(m.url())}
	for _, u := range server.LANURLs(m.boundAddr(), m.ips) {
		srv = append(srv, dim.Render("LAN ")+bold.Render(u))
	}
	srv = append(srv, dim.Render("uptime ")+m.uptime(), "", m.reloadButton())
	return panel("SERVER", "", srv, side, len(srv)+2) + "\n" + panel("ROOM", "", m.roomBody(), side, max(h-len(srv)-2, 2))
}

// clean drops control characters: names are client text and must not reach the terminal raw.
func clean(s string) string {
	if s = strings.Map(dropControl, s); s == "" {
		return dim.Render("—")
	}
	return s
}

func (m *model) roomBody() []string {
	s := m.sum
	var steps []string
	for _, st := range []struct{ status, label string }{
		{room.StatusMenu, "menu"}, {room.StatusLobby, "lobby"}, {room.StatusCreating, "create"}, {room.StatusPlaying, "play"},
	} {
		if st.status == s.Status {
			steps = append(steps, accent.Bold(true).Render(st.label))
		} else {
			steps = append(steps, dim.Render(st.label))
		}
	}
	dm := dim.Render("DM ") + clean(s.DM)
	if s.CoDM != "" {
		dm += dim.Render(" + ") + clean(s.CoDM)
	}
	lines := []string{
		accent.Bold(true).Render(strings.ToUpper(s.Status)),
		strings.Join(steps, dim.Render("▸")),
		clean(s.Campaign),
		dm,
		dim.Render("seats ") + bar(s.Players, s.Seats, 10) + fmt.Sprintf(" %d/%d", s.Players, s.Seats),
		dim.Render("ready ") + bar(s.Ready, s.Players, 10) + fmt.Sprintf(" %d/%d", s.Ready, s.Players),
		fmt.Sprintf("chars %d · scenes %d", s.Characters, s.Scenes),
	}
	if s.Scene != "" {
		lines = append(lines, dim.Render("scene ")+clean(s.Scene))
	}
	return lines
}

func (m *model) rightColumn(h int) string {
	n := side - 4 - 5 // sparkline cells, room for the " NNs" label
	conns := last(m.connHist, n)
	rej := fmt.Sprint(m.rejects)
	if m.rejects > 0 {
		rej = warn.Render(rej)
	}
	return panel("CONNECTIONS", "", []string{
		accent.Bold(true).Render(fmt.Sprint(m.clients)) + " online",
		accent.Render(spark(conns)) + dim.Render(fmt.Sprintf(" %ds", len(conns))),
		dim.Render(fmt.Sprintf("peak %d", peak(m.connHist))),
	}, side, 5) + "\n" + panel("ACTIVITY", "", []string{
		fmt.Sprintf("events %d · rej ", m.events) + rej,
		accent.Render(spark(last(m.evHist, n))) + dim.Render(" /10s"),
	}, side, 4) + "\n" + panel("SYSTEM", "", []string{
		fmt.Sprintf("heap %.1f MB", float64(m.heap)/(1<<20)),
		fmt.Sprintf("goroutines %d", m.goroutines),
	}, side, max(h-9, 2))
}

// statusLine is the one-line header of the narrow layout.
func (m *model) statusLine() string {
	link := bold.Render(m.url())
	if lan := server.LANURLs(m.boundAddr(), m.ips); len(lan) > 0 {
		link = dim.Render("LAN ") + bold.Render(lan[0])
	}
	parts := []string{m.statusText(), m.reloadButton(), link, fmt.Sprintf("%d online", m.clients),
		strings.ToUpper(m.sum.Status), "up " + m.uptime()}
	return ansi.Truncate(" "+strings.Join(parts, dim.Render(" · ")), max(m.w, 1), "…")
}

func bar(n, total, w int) string {
	f := 0
	if total > 0 {
		f = min(max(n, 0), total) * w / total
	}
	return accent.Render(strings.Repeat("█", f)) + dim.Render(strings.Repeat("░", w-f))
}

func spark(vals []int) string {
	const levels = "▁▂▃▄▅▆▇█"
	r := []rune(levels)
	hi := max(peak(vals), 1)
	var b strings.Builder
	for _, v := range vals {
		b.WriteRune(r[min(max(v, 0), hi)*(len(r)-1)/hi])
	}
	return b.String()
}

func peak(vals []int) int {
	p := 0
	for _, v := range vals {
		p = max(p, v)
	}
	return p
}

func last(vals []int, n int) []int { return vals[max(len(vals)-n, 0):] }

// push appends v to a ring of at most histLen samples.
func push(vals []int, v int) []int {
	vals = append(vals, v)
	return vals[max(len(vals)-histLen, 0):]
}

// hit reports whether (x, y) is on the first rendered occurrence of label, padded by px columns and py rows.
func hit(view, label string, x, y, px, py int) bool {
	for i, line := range strings.Split(ansi.Strip(view), "\n") {
		j := strings.Index(line, label)
		if j < 0 {
			continue
		}
		col := ansi.StringWidth(line[:j])
		return y >= i-py && y <= i+py && x >= col-px && x < col+ansi.StringWidth(label)+px
	}
	return false
}
