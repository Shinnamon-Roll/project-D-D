package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"
)

// Reload (docs/design/tui.md "v0.6 reload"): rebuild web + Go from the repo root, then re-exec the new binary.

const (
	buildTimeout   = 5 * time.Minute
	confirmTimeout = 5 * time.Second
)

// buildFunc runs the build in root and reports each output line (faked in tests).
type buildFunc func(ctx context.Context, root string, line func(string)) error

// findRoot locates the repo root: first the executable's repo (root/bin/tabletop), then walking up from wd.
// The exe check comes first so a stray cwd cannot pick the build dir; under `go run` the exe is a temp build.
func findRoot(exe, wd string) (root, reason string) {
	if exe != "" {
		if real, err := filepath.EvalSymlinks(exe); err == nil && isRoot(filepath.Dir(filepath.Dir(real))) {
			return filepath.Dir(filepath.Dir(real)), ""
		}
	}
	if wd == "" {
		return "", "no working directory"
	}
	for d := wd; ; d = filepath.Dir(d) {
		if isRoot(d) {
			return d, ""
		}
		if filepath.Dir(d) == d {
			return "", "no projectdnd repo (Makefile, go.mod, cmd/tabletop) at the binary or above " + wd
		}
	}
}

// isRoot: d has a Makefile, cmd/tabletop and a go.mod declaring module projectdnd.
func isRoot(d string) bool {
	mk, err1 := os.Stat(filepath.Join(d, "Makefile"))
	src, err2 := os.Stat(filepath.Join(d, "cmd", "tabletop"))
	mod, err3 := os.ReadFile(filepath.Join(d, "go.mod"))
	if err1 != nil || err2 != nil || err3 != nil || !mk.Mode().IsRegular() || !src.IsDir() {
		return false
	}
	for l := range strings.Lines(string(mod)) {
		if strings.TrimSpace(l) == "module projectdnd" {
			return true
		}
	}
	return false
}

// makeBuild mirrors `make run`: web build, then `make tabletop` (go build with the version ldflags).
func makeBuild(ctx context.Context, root string, line func(string)) error {
	steps := []struct {
		dir  string
		args []string
	}{
		{"web", []string{"npm", "run", "build"}},
		{".", []string{"make", "tabletop"}},
	}
	for _, s := range steps {
		name := strings.Join(s.args, " ")
		line("$ " + name)
		cmd := exec.CommandContext(ctx, s.args[0], s.args[1:]...)
		cmd.Dir = filepath.Join(root, s.dir)
		// own process group: cancel/timeout kills npm/make and every child, not just the direct process
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
		cmd.WaitDelay = 5 * time.Second
		w := &lineWriter{line: line}
		cmd.Stdout, cmd.Stderr = w, w // same writer: exec serializes Write calls
		err := cmd.Run()
		w.flush()
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

// lineWriter splits process output into lines (\n or \r, for progress bars).
type lineWriter struct {
	buf  []byte
	line func(string)
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexAny(w.buf, "\r\n")
		if i < 0 {
			break
		}
		if i > 0 {
			w.line(string(w.buf[:i]))
		}
		w.buf = w.buf[i+1:]
	}
	if len(w.buf) > 4096 { // no newline in sight: emit rather than grow forever
		w.flush()
	}
	return len(p), nil
}

func (w *lineWriter) flush() {
	if len(w.buf) > 0 {
		w.line(string(w.buf))
		w.buf = nil
	}
}

var execve = syscall.Exec // replaced in tests

// reexec replaces this process with the freshly built root/bin/tabletop, same args + -autostart.
// The rebuilt binary is used, not os.Executable(): under `go run` that is a stale temp build.
func reexec(root string, args []string) error {
	bin := filepath.Join(root, "bin", "tabletop")
	argv := slices.Concat([]string{bin}, args)
	if !slices.Contains(args, "-autostart") {
		argv = append(argv, "-autostart")
	}
	return execve(bin, argv, os.Environ())
}
