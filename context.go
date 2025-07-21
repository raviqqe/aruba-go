package aruba

import (
	"io"
	"os/exec"
)

type contextKey struct{}

type arubaContext struct {
	Command   *exec.Cmd
	Directory string
	Stdin     io.WriteCloser
}
