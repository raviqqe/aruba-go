package aruba

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

func parseString(s string) (string, error) {
	return strconv.Unquote(s)
}

func parseDocString(s string) string {
	return strings.TrimRight(s, "\n")
}

func matchesExactly(s, t string) bool {
	return s == t || strings.TrimSpace(s) == t
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "aruba-*")

	return contextWithWorld(ctx, world{Directory: d}), err
}

func createFile(ctx context.Context, p string, docString *godog.DocString) error {
	return os.WriteFile(
		path.Join(contextWorld(ctx).Directory, p),
		[]byte(parseDocString(docString.Content)+"\n"),
		0o600,
	)
}

func createDirectory(ctx context.Context, p string) error {
	return os.Mkdir(path.Join(contextWorld(ctx).Directory, p), 0o700)
}

func runCommand(ctx context.Context, successfully, command, interactively string) (context.Context, error) {
	command, err := parseString(command)
	if err != nil {
		return ctx, err
	}

	ss := strings.Split(command, " ")
	c := exec.Command(ss[0], ss[1:]...)
	w := contextWorld(ctx)
	c.Dir = w.Directory
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout
	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr
	w.Command = c
	ctx = contextWithWorld(ctx, w)

	in, err := c.StdinPipe()
	if err != nil {
		return ctx, err
	}

	w.Stdin = in
	ctx = contextWithWorld(ctx, w)

	err = c.Start()
	if err != nil {
		return ctx, err
	}

	if interactively == "" {
		if err := c.Wait(); successfully != "" {
			return ctx, err
		}
	}

	return ctx, nil
}

func exitStatus(ctx context.Context, not string, code int) error {
	c := contextWorld(ctx).Command
	_ = c.Wait()

	if c := c.ProcessState.ExitCode(); (c == code) != (not == "") {
		return fmt.Errorf("expected exit code %d%s to be %d", c, not, code)
	}

	return nil
}

func stdin(ctx context.Context, p string) error {
	f, err := os.Open(path.Join(contextWorld(ctx).Directory, p))
	if err != nil {
		return err
	}

	go func() {
		w := contextWorld(ctx).Stdin
		_, _ = io.Copy(w, f)
		_ = w.Close()
	}()

	return err
}

func stdout(ctx context.Context, stdout, not, exactly, pattern string) error {
	c := contextWorld(ctx).Command
	_ = c.Wait()
	out := c.Stdout

	if stdout == "stderr" {
		out = c.Stderr
	}

	s := out.(*bytes.Buffer).String()

	if exactly == "" && strings.Contains(s, pattern) != (not == "") ||
		exactly != "" && matchesExactly(s, pattern) != (not == "") {
		return fmt.Errorf("expected %s %q%s to contain%s %q", stdout, s, not, exactly, pattern)
	}

	return nil
}

func fileContains(ctx context.Context, p, not, exactly, pattern string) error {
	bs, err := os.ReadFile(path.Join(contextWorld(ctx).Directory, p))
	if err != nil {
		return err
	}

	s := string(bs)
	ok := strings.Contains(s, pattern)

	fmt.Printf("%q %q %q", exactly, s, pattern)
	if exactly != "" {
		ok = matchesExactly(s, pattern)
	}

	if ok != (not == "") {
		return fmt.Errorf("expected file %q%s to contain%s %q", p, not, exactly, pattern)
	}

	return nil
}

func fileExists(ctx context.Context, ty, p, not string) error {
	if i, err := os.Stat(path.Join(contextWorld(ctx).Directory, p)); (err == nil && i.IsDir() == (ty == "directory")) != (not == "") {
		return fmt.Errorf("%s %q should%s exist", ty, p, not)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(before)
	ctx.Step(`^a file named "(.+)" with:$`, createFile)
	ctx.Step(`^a directory named "(.+)"$`, createDirectory)
	ctx.Step("^I( successfully)? run (`.*`)( interactively)?$", runCommand)
	ctx.Step(`^the exit status should( not)? be (\d+)$`, exitStatus)
	ctx.Step(
		`^the (std(?:out|err)) should( not)? contain( exactly)? (".*")$`,
		func(ctx context.Context, port, not, exactly, pattern string) error {
			pattern, err := parseString(pattern)
			if err != nil {
				return err
			}

			return stdout(ctx, port, not, exactly, pattern)
		},
	)
	ctx.Step(
		`^the (std(?:out|err)) should( not)? contain( exactly)?:$`,
		func(ctx context.Context, port, not, exactly string, docString *godog.DocString) error {
			return stdout(ctx, port, not, exactly, parseDocString(docString.Content))
		},
	)
	ctx.Step(`^a file named "(.*)" should( not)? contain (".*")$`, func(ctx context.Context, p, not, pattern string) error {
		pattern, err := parseString(pattern)
		if err != nil {
			return err
		}

		return fileContains(ctx, p, not, "", pattern)
	})
	ctx.Step(`^a file named "(.+)" should( not)? contain( exactly)?:$`, func(ctx context.Context, p, not, exactly string, docString *godog.DocString) error {
		return fileContains(ctx, p, not, exactly, parseDocString(docString.Content))
	})
	ctx.Step(`^I pipe in the file(?: named)? "(.*)"$`, stdin)
	ctx.Step(`^(a|the) (directory|file)(?: named)? "(.*)" should( not)? exist$`, fileExists)
}
