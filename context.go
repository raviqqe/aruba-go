package aruba

import (
	"context"
	"io"
	"os/exec"
)

type worldKey struct{}

type world struct {
	Command   *exec.Cmd
	Directory string
	Stdin     io.WriteCloser
}

func contextWorld(ctx context.Context) world {
	return ctx.Value(worldKey{}).(world)
}

func contextWithWorld(ctx context.Context, w world) context.Context {
	return context.WithValue(ctx, worldKey{}, w)
}
