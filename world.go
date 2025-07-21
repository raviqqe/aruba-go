package aruba

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

type worldKey struct{}

type world struct {
	Command          *exec.Cmd
	RootDirectory    string
	CurrentDirectory string
	Stdin            io.WriteCloser
	Stdout           *bytes.Buffer
	Stderr           *bytes.Buffer
}

func newWorld(d string) world {
	return world{
		RootDirectory:    d,
		CurrentDirectory: d,
		Stdout:           bytes.NewBuffer(nil),
		Stderr:           bytes.NewBuffer(nil),
	}
}

func contextWorld(ctx context.Context) world {
	return ctx.Value(worldKey{}).(world)
}

func contextWithWorld(ctx context.Context, w world) context.Context {
	return context.WithValue(ctx, worldKey{}, w)
}
