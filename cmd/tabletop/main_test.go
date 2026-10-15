package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

type fakeServer struct{ shutdowns int }

func (f *fakeServer) Addr() string                   { return "127.0.0.1:9999" }
func (f *fakeServer) Clients() int                   { return 3 }
func (f *fakeServer) Shutdown(context.Context) error { f.shutdowns++; return nil }

func key(s string) tea.KeyPressMsg {
	switch s {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "space":
		return tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}
	case "ctrl+c":
		return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	}
	return tea.KeyPressMsg{Code: rune(s[0]), Text: s}
}

func newTestModel(start func() (runner, error)) *model {
	m := newModel("127.0.0.1:8080", "0.3.0", start)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	return m
}

func view(m *model) string { return ansi.Strip(m.View().Content) }

// runCmd runs cmd and returns the first non-batch message it produces.
func runCmd(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		t.Fatal("nil cmd")
	}
	msg := cmd()
	if b, ok := msg.(tea.BatchMsg); ok {
		for _, c := range b {
			if c == nil {
				continue
			}
			if m := c(); m != nil {
				if _, tick := m.(spinner.TickMsg); !tick {
					return m
				}
			}
		}
	}
	return msg
}

func TestWelcomeStartFlow(t *testing.T) {
	for _, k := range []string{"enter", "space"} {
		t.Run(k, func(t *testing.T) {
			srv := &fakeServer{}
			m := newTestModel(func() (runner, error) { return srv, nil })
			if !strings.Contains(view(m), "เริ่มเซิร์ฟเวอร์") || !strings.Contains(view(m), "v0.3.0") {
				t.Fatalf("welcome view:\n%s", view(m))
			}
			_, cmd := m.Update(key(k))
			if m.screen != starting || !strings.Contains(view(m), "กำลังเริ่ม") {
				t.Fatalf("screen %v after %s", m.screen, k)
			}
			m.Update(runCmd(t, cmd))
			if m.screen != running || !strings.Contains(view(m), "RUNNING") || !strings.Contains(view(m), "clients 3") {
				t.Fatalf("running view:\n%s", view(m))
			}
		})
	}
}

func TestStartFailureShowsRetry(t *testing.T) {
	calls := 0
	m := newTestModel(func() (runner, error) {
		calls++
		return nil, errors.New("address already in use")
	})
	_, cmd := m.Update(key("enter"))
	m.Update(runCmd(t, cmd))
	v := view(m)
	if m.screen != welcome || !strings.Contains(v, "address already in use") || !strings.Contains(v, "ลองใหม่") {
		t.Fatalf("after failure:\n%s", v)
	}
	_, cmd = m.Update(key("enter"))
	runCmd(t, cmd)
	if calls != 2 {
		t.Fatalf("retry did not call start again: %d", calls)
	}
}

func TestMouseClickOnButton(t *testing.T) {
	m := newTestModel(func() (runner, error) { return &fakeServer{}, nil })
	lines := strings.Split(view(m), "\n")
	y, x := -1, -1
	for i, l := range lines {
		if j := strings.Index(l, "เริ่มเซิร์ฟเวอร์"); j >= 0 {
			y, x = i, ansi.StringWidth(l[:j])
		}
	}
	if y < 0 {
		t.Fatal("button not rendered")
	}
	m.Update(tea.MouseClickMsg{X: 0, Y: 0, Button: tea.MouseLeft})
	if m.screen != welcome {
		t.Fatal("click outside the button started the server")
	}
	m.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
	if m.screen != starting {
		t.Fatal("click on the button did not start")
	}
}

func TestLogRingBuffer(t *testing.T) {
	m := newTestModel(nil)
	for i := range maxLines + 5 {
		m.Update(logMsg{time.Now(), slog.LevelInfo, fmt.Sprintf("line %d", i)})
	}
	first, last := ansi.Strip(m.lines[0]), ansi.Strip(m.lines[maxLines-1])
	if len(m.lines) != maxLines || !strings.HasSuffix(first, "INFO  line 5") || !strings.HasSuffix(last, "line 1004") {
		t.Fatalf("len %d first %q last %q", len(m.lines), first, last)
	}
	m.screen, m.srv = running, &fakeServer{}
	m.Update(key("c"))
	if len(m.lines) != 0 {
		t.Fatal("c did not clear")
	}
}

func TestQuitShutsDownServer(t *testing.T) {
	for _, k := range []string{"q", "ctrl+c"} {
		t.Run(k, func(t *testing.T) {
			srv := &fakeServer{}
			m := newTestModel(func() (runner, error) { return srv, nil })
			m.Update(startedMsg{srv: srv})
			_, cmd := m.Update(key(k))
			if !m.stopping || !strings.Contains(view(m), "STOPPING") {
				t.Fatal("not stopping")
			}
			if _, again := m.Update(key(k)); again != nil {
				t.Fatal("second q started another shutdown")
			}
			msg := runCmd(t, cmd)
			if srv.shutdowns != 1 {
				t.Fatalf("shutdowns %d", srv.shutdowns)
			}
			_, cmd = m.Update(msg)
			if _, ok := runCmd(t, cmd).(tea.QuitMsg); !ok || !m.stopped {
				t.Fatal("did not quit after shutdown")
			}
		})
	}
	m := newTestModel(nil)
	_, cmd := m.Update(key("q"))
	if _, ok := runCmd(t, cmd).(tea.QuitMsg); !ok {
		t.Fatal("q on welcome did not quit")
	}
}

func TestFollowToggle(t *testing.T) {
	m := newTestModel(nil)
	m.screen, m.srv = running, &fakeServer{}
	for i := range 100 {
		m.Update(logMsg{time.Now(), slog.LevelInfo, fmt.Sprint(i)})
	}
	if !m.vp.AtBottom() {
		t.Fatal("follow did not stick to bottom")
	}
	m.Update(key("f"))
	m.Update(logMsg{time.Now(), slog.LevelInfo, "x"})
	if m.follow || m.vp.AtBottom() {
		t.Fatal("follow off still scrolled")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyPgUp})
	if m.follow {
		t.Fatal("scroll up kept follow")
	}
}

func TestTeaHandler(t *testing.T) {
	ch := make(chan logMsg, 1)
	log := slog.New(&teaHandler{ch: ch}).With("srv", "a").WithGroup("g")
	log.Warn("event rejected", "type", "token.move", "err", "you can't\x1b[2J")
	got := <-ch
	if got.level != slog.LevelWarn {
		t.Fatal(got.level)
	}
	want := `event rejected srv=a g.type=token.move g.err="you can't\x1b[2J"`
	if got.text != want {
		t.Fatalf("got  %s\nwant %s", got.text, want)
	}

	// Full channel: Handle must drop, never block.
	ch <- got
	done := make(chan struct{})
	go func() { log.Info("dropped"); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Handle blocked on a full channel")
	}
}

func TestShutdownIfRunning(t *testing.T) {
	tests := []struct {
		name          string
		srv           *fakeServer
		stopped       bool
		wantShutdowns int
	}{
		{"never started", nil, false, 0},
		{"running, q path skipped", &fakeServer{}, false, 1},
		{"already stopped via q", &fakeServer{}, true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(nil)
			if tt.srv != nil {
				m.srv = tt.srv
			}
			m.stopped = tt.stopped
			m.shutdownIfRunning()
			m.shutdownIfRunning() // idempotent
			if tt.srv != nil && tt.srv.shutdowns != tt.wantShutdowns {
				t.Fatalf("shutdowns = %d, want %d", tt.srv.shutdowns, tt.wantShutdowns)
			}
		})
	}
}
