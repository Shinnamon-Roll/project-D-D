package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"projectdnd/internal/room"
)

// runningModel is a started dashboard of w×h with a fake summary, some logs and samples.
func runningModel(t *testing.T, w, h int, srv *fakeServer) *model {
	t.Helper()
	m := newModel("0.0.0.0:8080", "0.6.0", func() (runner, error) { return srv, nil })
	m.ips = nil
	m.root = "/repo"
	m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	m.Update(startedMsg{srv: srv})
	return m
}

func fakeSummary() room.Summary {
	return room.Summary{Status: room.StatusLobby, Campaign: "Starter Set", DM: "Dee Em", CoDM: "Cody",
		Players: 3, Seats: 6, Ready: 2, Characters: 0, Scenes: 1, Scene: "Beach"}
}

func TestDashboardLayouts(t *testing.T) {
	tests := []struct {
		w, h       int
		want, none []string
	}{
		{160, 45, []string{"SERVER", "ROOM", "LOGS", "CONNECTIONS", "ACTIVITY", "SYSTEM", "3 online", "seats", "Dee Em + Cody", reloadLabel}, nil},
		{110, 35, []string{"SERVER", "ROOM", "LOGS", "3 online · ev 1 · rej 1", reloadLabel}, []string{"CONNECTIONS", "SYSTEM"}},
		{80, 24, []string{"RUNNING", "LOGS", "3 online", "LOBBY", reloadLabel}, []string{"SERVER", "ROOM", "CONNECTIONS"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%dx%d", tt.w, tt.h), func(t *testing.T) {
			m := runningModel(t, tt.w, tt.h, &fakeServer{clients: 3, sum: fakeSummary()})
			m.Update(logMsg{time.Now(), slog.LevelInfo, "event token.move by=Mira"})
			m.Update(logMsg{time.Now(), slog.LevelWarn, `event rejected type=token.move by=Bo err="you can only move your own token"`})
			m.Update(logMsg{time.Now(), slog.LevelInfo, strings.Repeat("long line ", 40)})
			for range 25 {
				m.Update(tickMsg(time.Now()))
			}
			v := view(m)
			lines := strings.Split(v, "\n")
			if len(lines) > tt.h {
				t.Errorf("%d lines > height %d", len(lines), tt.h)
			}
			for i, l := range strings.Split(m.View().Content, "\n") {
				if w := ansi.StringWidth(l); w > tt.w {
					t.Errorf("line %d width %d > %d: %q", i, w, tt.w, ansi.Strip(l))
				}
			}
			for _, s := range tt.want {
				if !strings.Contains(v, s) {
					t.Errorf("missing %q", s)
				}
			}
			for _, s := range tt.none {
				if strings.Contains(v, s) {
					t.Errorf("unexpected %q", s)
				}
			}
			if t.Failed() {
				t.Logf("view:\n%s", v)
			}
		})
	}
	// Tiny and unsized terminals must not panic.
	for _, sz := range [][2]int{{0, 0}, {1, 1}, {10, 3}, {90, 5}, {130, 4}} {
		runningModel(t, sz[0], sz[1], &fakeServer{sum: fakeSummary()}).View()
	}
}

func TestDashboardCountsAndSanitizes(t *testing.T) {
	sum := fakeSummary()
	sum.DM, sum.Scene = "Ev\x1b[2Jil", "Bea\u009b31mch"
	m := runningModel(t, 160, 45, &fakeServer{sum: sum})
	m.Update(logMsg{time.Now(), slog.LevelInfo, "event seat.add by=Dee"})
	m.Update(logMsg{time.Now(), slog.LevelWarn, "event rejected type=x by=- err=no"})
	m.Update(logMsg{time.Now(), slog.LevelError, "boom"})
	if m.events != 1 || m.rejects != 1 || m.warns != 2 {
		t.Fatalf("events %d rejects %d warns %d", m.events, m.rejects, m.warns)
	}
	raw := m.View().Content
	if strings.Contains(raw, "\x1b[2J") || strings.Contains(raw, "\u009b") || !strings.Contains(view(m), "Ev[2Jil") {
		t.Fatal("client text reached the terminal unsanitized")
	}
	for range histLen + 5 {
		m.Update(tickMsg(time.Now()))
	}
	if len(m.connHist) != histLen || len(m.evHist) != (histLen+6)/10 || m.evHist[0] != 2 {
		t.Fatalf("connHist %d evHist %v", len(m.connHist), m.evHist)
	}
}

// fakeBuild records calls and returns err.
type fakeBuild struct {
	calls int
	err   error
}

func (f *fakeBuild) run(_ context.Context, root string, line func(string)) error {
	f.calls++
	line("\x1b[32mvite built\x1b[0m in " + root)
	return f.err
}

func TestReloadConfirm(t *testing.T) {
	tests := []struct {
		name      string
		clients   int
		root      string
		presses   int
		expired   bool // the confirm window ran out before the 2nd press
		wantBuild bool
	}{
		{"no clients: builds at once", 0, "/repo", 1, false, true},
		{"clients: first r only asks", 2, "/repo", 1, false, false},
		{"clients: second r confirms", 2, "/repo", 2, false, true},
		{"clients: expired confirm asks again", 2, "/repo", 2, true, false},
		{"no repo: disabled", 0, "", 2, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := runningModel(t, 160, 45, &fakeServer{clients: tt.clients})
			fb := &fakeBuild{}
			m.root, m.noReload, m.build = tt.root, "no Makefile", fb.run
			var cmd tea.Cmd
			for i := range tt.presses {
				if i > 0 && tt.expired {
					m.confirm = time.Now().Add(-time.Millisecond)
				}
				_, cmd = m.Update(key("r"))
			}
			if m.building != tt.wantBuild || (cmd != nil) != tt.wantBuild {
				t.Fatalf("building %v cmd %v, want %v", m.building, cmd != nil, tt.wantBuild)
			}
			v := view(m)
			if tt.root == "" && (!strings.Contains(v, "reload disabled: no Makefile") || strings.Contains(v, "RELOAD off:")) {
				t.Fatalf("no reason shown:\n%s", v)
			}
			if !tt.wantBuild && tt.root != "" && !strings.Contains(v, "กด r อีกครั้ง") {
				t.Fatalf("no confirm prompt:\n%s", v)
			}
			if tt.wantBuild {
				if _, again := m.Update(key("r")); again != nil {
					t.Fatal("r while building started another build")
				}
			}
		})
	}
}

func TestReloadBuildResult(t *testing.T) {
	for _, fail := range []bool{true, false} {
		t.Run(fmt.Sprint("fail=", fail), func(t *testing.T) {
			srv := &fakeServer{}
			m := runningModel(t, 160, 45, srv)
			ch := make(chan logMsg, 8)
			m.log = slog.New(&teaHandler{ch: ch})
			fb := &fakeBuild{}
			if fail {
				fb.err = errors.New("exit status 2")
			}
			m.build = fb.run
			_, cmd := m.Update(key("r"))
			if !strings.Contains(view(m), "BUILDING") {
				t.Fatal("no BUILDING status")
			}
			done := runCmd(t, cmd)
			if got := (<-ch).text; got != "build │ vite built in /repo" {
				t.Fatalf("build line %q", got)
			}
			_, cmd = m.Update(done)
			if fail {
				if cmd != nil || srv.shutdowns != 0 || m.reload || m.building || m.screen != running ||
					!strings.Contains(view(m), "build failed, server keeps running") {
					t.Fatalf("failed build: cmd %v shutdowns %d reload %v\n%s", cmd != nil, srv.shutdowns, m.reload, view(m))
				}
				return
			}
			if !m.reload || !m.stopping {
				t.Fatal("successful build did not stop for reload")
			}
			_, cmd = m.Update(runCmd(t, cmd))
			if _, ok := runCmd(t, cmd).(tea.QuitMsg); !ok || srv.shutdowns != 1 {
				t.Fatalf("did not shut down + quit: shutdowns %d", srv.shutdowns)
			}
		})
	}
}

func TestQuitDuringBuildNoReload(t *testing.T) {
	srv := &fakeServer{}
	m := runningModel(t, 160, 45, srv)
	m.build = (&fakeBuild{}).run
	_, build := m.Update(key("r"))
	m.Update(key("q"))
	m.Update(runCmd(t, build))
	if m.reload {
		t.Fatal("build finishing after q still reloaded")
	}
}

func TestReloadClick(t *testing.T) {
	for _, w := range []int{160, 110, 80} {
		m := runningModel(t, w, 30, &fakeServer{})
		fb := &fakeBuild{}
		m.build = fb.run
		y, x := -1, -1
		for i, l := range strings.Split(view(m), "\n") {
			if j := strings.Index(l, reloadLabel); j >= 0 {
				y, x = i, ansi.StringWidth(l[:j])
				break
			}
		}
		m.Update(tea.MouseClickMsg{X: x, Y: y + 1, Button: tea.MouseLeft})
		if m.building {
			t.Fatalf("w=%d: click below the button reloaded", w)
		}
		m.Update(tea.MouseClickMsg{X: x + 2, Y: y, Button: tea.MouseLeft})
		if !m.building {
			t.Fatalf("w=%d: click on %s at %d,%d did nothing", w, reloadLabel, x, y)
		}
	}
}

func TestAutostart(t *testing.T) {
	m := newTestModel(func() (runner, error) { return &fakeServer{}, nil })
	if m.Init() != nil {
		t.Fatal("Init without -autostart did something")
	}
	m.autostart = true
	msg := runCmd(t, m.Init())
	m.Update(msg)
	if m.screen != running {
		t.Fatalf("autostart screen %v", m.screen)
	}
}

func TestReexec(t *testing.T) {
	var gotBin string
	var gotArgv []string
	defer func(old func(string, []string, []string) error) { execve = old }(execve)
	execve = func(bin string, argv, _ []string) error { gotBin, gotArgv = bin, argv; return errors.New("faked") }
	for _, args := range [][]string{{"-addr", "127.0.0.1:9000"}, {"-addr", "127.0.0.1:9000", "-autostart"}} {
		if err := reexec("/repo", args); err == nil {
			t.Fatal("exec error swallowed")
		}
		want := []string{"/repo/bin/tabletop", "-addr", "127.0.0.1:9000", "-autostart"}
		if gotBin != want[0] || !slices.Equal(gotArgv, want) {
			t.Fatalf("args %v: exec %s %v", args, gotBin, gotArgv)
		}
	}
}

func TestFindRoot(t *testing.T) {
	mkRepo := func(module string, makefile bool) string {
		root := t.TempDir()
		for _, d := range []string{filepath.Join(root, "web", "src"), filepath.Join(root, "cmd", "tabletop"), filepath.Join(root, "bin")} {
			if err := os.MkdirAll(d, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		files := map[string]string{"go.mod": "module " + module + "\n\ngo 1.26\n", "bin/tabletop": ""}
		if makefile {
			files["Makefile"] = ""
		}
		for f, body := range files {
			if err := os.WriteFile(filepath.Join(root, f), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	repo, other := mkRepo("projectdnd", true), mkRepo("projectdnd", true)
	noMake, wrongMod := mkRepo("projectdnd", false), mkRepo("projectdnd2", true)
	link := filepath.Join(t.TempDir(), "tabletop")
	if err := os.Symlink(filepath.Join(repo, "bin", "tabletop"), link); err != nil {
		t.Fatal(err)
	}
	tmpExe := filepath.Join(t.TempDir(), "exe") // go run: temp binary outside any repo
	tests := []struct {
		name, exe, wd, want string
	}{
		{"exe wins over cwd", filepath.Join(repo, "bin", "tabletop"), other, repo},
		{"exe via symlink", link, "", repo},
		{"go run: cwd walk-up", tmpExe, filepath.Join(repo, "web", "src"), repo},
		{"no exe: cwd at root", "", repo, repo},
		{"no Makefile", filepath.Join(noMake, "bin", "tabletop"), filepath.Join(noMake, "web"), ""},
		{"other module", filepath.Join(wrongMod, "bin", "tabletop"), wrongMod, ""},
		{"nothing", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, why := findRoot(tt.exe, tt.wd)
			if tt.want != "" { // macOS TempDir is under the /var -> /private/var symlink
				tt.want, _ = filepath.EvalSymlinks(tt.want)
				got, _ = filepath.EvalSymlinks(got)
			}
			if got != tt.want || (got == "") != (why != "") {
				t.Fatalf("findRoot(%q, %q) = %q %q, want %q", tt.exe, tt.wd, got, why, tt.want)
			}
		})
	}
}

func TestLineWriter(t *testing.T) {
	var got []string
	w := &lineWriter{line: func(s string) { got = append(got, s) }}
	for _, p := range []string{"a\nb", "c\r\n", "\rd"} {
		w.Write([]byte(p))
	}
	w.flush()
	if !slices.Equal(got, []string{"a", "bc", "d"}) {
		t.Fatalf("%q", got)
	}
}
