package aruba

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type worldKey struct{}

type world struct {
	commands         []*exec.Cmd
	RootDirectory    string
	CurrentDirectory string
	Stdin            io.WriteCloser
	Environment      []string
	StartupWaitTime  time.Duration
}

func newWorld(d string) world {
	return world{
		RootDirectory:    d,
		CurrentDirectory: d,
		Environment:      os.Environ(),
	}
}

func (w world) AddCommand(c *exec.Cmd) world {
	w.commands = append(w.commands, c)
	return w
}

func (w world) FindCommand(s string) *exec.Cmd {
	for _, c := range w.commands {
		if s == strings.Join(c.Args, " ") {
			return c
		}
	}

	return nil
}

func (w world) LastCommand() *exec.Cmd {
	return w.commands[len(w.commands)-1]
}

func (w world) Stop() {
	for _, c := range w.commands {
		if c.Process != nil && c.ProcessState == nil {
			_ = c.Process.Kill()
			_ = c.Wait()
		}
	}
}

func (w world) Stdout() string {
	return w.stdout(func(c *exec.Cmd) io.Writer {
		return c.Stdout
	})
}

func (w world) Stderr() string {
	return w.stdout(func(c *exec.Cmd) io.Writer {
		return c.Stderr
	})
}

func (w world) Output() string {
	bs := []byte(nil)

	for _, c := range w.commands {
		_ = c.Wait()

		for _, b := range []*bytes.Buffer{c.Stdout.(*bytes.Buffer), c.Stderr.(*bytes.Buffer)} {
			bs = append(bs, b.Bytes()...)
		}
	}

	return string(bs)
}

func (w world) path(p string) (string, error) {
	if filepath.IsAbs(p) || filepath.VolumeName(p) != "" {
		return "", fmt.Errorf("path %q must be relative to the working directory", p)
	}

	q := filepath.Join(w.CurrentDirectory, p)

	if d, err := filepath.Rel(w.RootDirectory, q); err != nil {
		return "", err
	} else if strings.HasPrefix(d, "..") {
		return "", fmt.Errorf("path %q is outside the working directory", p)
	}

	return q, nil
}

func (w world) stdout(f func(*exec.Cmd) io.Writer) string {
	bs := []byte(nil)

	for _, c := range w.commands {
		_ = c.Wait()
		bs = append(bs, f(c).(*bytes.Buffer).Bytes()...)
	}

	return string(bs)
}

func contextWorld(ctx context.Context) world {
	return ctx.Value(worldKey{}).(world)
}

func contextWithWorld(ctx context.Context, w world) context.Context {
	return context.WithValue(ctx, worldKey{}, w)
}
