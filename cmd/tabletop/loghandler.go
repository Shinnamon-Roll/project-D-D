package main

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// logMsg is one slog record, formatted, on its way to the TUI.
type logMsg struct {
	time  time.Time
	level slog.Level
	text  string // message + key=value attrs
}

// teaHandler forwards slog records to a buffered channel that a goroutine drains into
// tea.Program.Send. Handle never blocks: when the TUI falls behind, records are dropped.
type teaHandler struct {
	ch     chan<- logMsg
	attrs  string // pre-formatted WithAttrs
	prefix string // group prefix, "g." per WithGroup
}

func (h *teaHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *teaHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	// Messages can carry outside text (http.Server.ErrorLog), so strip control chars like attr values.
	b.WriteString(strings.Map(dropControl, r.Message))
	b.WriteString(h.attrs)
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(&b, h.prefix, a)
		return true
	})
	select {
	case h.ch <- logMsg{r.Time, r.Level, b.String()}:
	default: // TUI too slow: drop rather than block the server
	}
	return nil
}

func (h *teaHandler) WithAttrs(as []slog.Attr) slog.Handler {
	var b strings.Builder
	b.WriteString(h.attrs)
	for _, a := range as {
		writeAttr(&b, h.prefix, a)
	}
	return &teaHandler{h.ch, b.String(), h.prefix}
}

func (h *teaHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &teaHandler{h.ch, h.attrs, h.prefix + name + "."}
}

func writeAttr(b *strings.Builder, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		p := prefix
		if a.Key != "" {
			p += a.Key + "."
		}
		for _, g := range a.Value.Group() {
			writeAttr(b, p, g)
		}
		return
	}
	v := a.Value.String()
	// Quote spaces and control chars: values can carry client text (names, errors) and must not
	// inject terminal escape sequences into the TUI.
	if v == "" || strings.IndexFunc(v, func(r rune) bool { return r == ' ' || r == '"' || r == '=' || !unicode.IsPrint(r) }) >= 0 {
		v = strconv.Quote(v)
	}
	b.WriteString(" " + prefix + a.Key + "=" + v)
}

func dropControl(r rune) rune {
	if unicode.IsControl(r) {
		return -1
	}
	return r
}
