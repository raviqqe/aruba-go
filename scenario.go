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

func createFile(ctx context.Context, p, s string) error {
	return os.WriteFile(path.Join(contextWorld(ctx).Directory, p), []byte(s), 0o600)
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
	c.Stdout = bytes.NewBuffer(nil)
	c.Stderr = bytes.NewBuffer(nil)
	w.Command = c
	ctx = contextWithWorld(ctx, w)

	w.Stdin, err = c.StdinPipe()
	if err != nil {
		return ctx, err
	}

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

func runScript(ctx context.Context, s *godog.DocString) (context.Context, error) {
	const scriptPath = "script"

	err := createFile(ctx, scriptPath, s.Content)
	if err != nil {
		return ctx, err
	}

	return runCommand(ctx, "", strconv.Quote("sh "+scriptPath), "")
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
	w := contextWorld(ctx)
	f, err := os.Open(path.Join(w.Directory, p))
	if err != nil {
		return err
	}

	// TODO Figure out why we need to ignore errors...
	_, _ = io.Copy(w.Stdin, f)
	_ = w.Stdin.Close()

	return nil
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

func setEnvVar(ctx context.Context, k, v string) error {
	return os.Setenv(k, v)
}

func changeDirectory(ctx context.Context, p string) (context.Context, error) {
	w := contextWorld(ctx)

	return contextWithWorld(ctx, w), nil
}

// [InitializeScenario] initializes a scenario.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(before)
	ctx.Step(`^a file named "(.+)" with (".*")$`, func(ctx context.Context, p, s string) error {
		s, err := parseString(s)
		if err != nil {
			return err
		}
		return createFile(ctx, p, s)
	})
	ctx.Step(`^a file named "(.+)" with:$`, func(ctx context.Context, p string, s *godog.DocString) error {
		return createFile(ctx, p, parseDocString(s.Content)+"\n")
	})
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
	ctx.Step(`^(?:a|the) (directory|file)(?: named)? "(.*)" should( not)? exist$`, fileExists)
	ctx.Step(`^I set the environment variable "(.*)" to "(.*)"$`, setEnvVar)
	ctx.Step(`^I run the following script:$`, runScript)
	ctx.Step(`^I cd to "(.*)"$`, changeDirectory)
}
