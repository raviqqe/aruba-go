package aruba

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"
)

type worldKey struct{}

type world struct {
	commands         []*exec.Cmd
	RootDirectory    string
	CurrentDirectory string
	Stdin            io.WriteCloser
}

func newWorld(d string) world {
	return world{
		RootDirectory:    d,
		CurrentDirectory: d,
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
