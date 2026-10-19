package aruba

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

var stringUnescaper = strings.NewReplacer(`\"`, `"`, `\'`, `'`)

// Strings unescape only quotes like the string parameter type of Cucumber expressions.
func unescapeString(s string) string {
	return stringUnescaper.Replace(s)
}

func trimTrailingNewlines(s string) string {
	return strings.TrimRight(s, "\n")
}

func matchesExactly(s, t string) bool {
	return s == t || trimTrailingNewlines(s) == trimTrailingNewlines(t)
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "aruba-*")

	return contextWithWorld(ctx, newWorld(d)), err
}

func after(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	w := contextWorld(ctx)
	w.Stop()

	return ctx, os.RemoveAll(w.RootDirectory)
}

func createFile(ctx context.Context, p, s string) error {
	p, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}

	return os.WriteFile(p, []byte(s), 0o600)
}

func createFileWithMode(ctx context.Context, p, s string, mode os.FileMode) error {
	if err := createFile(ctx, p, s); err != nil {
		return err
	}

	q, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	}

	return os.Chmod(q, mode)
}

func parseFileMode(s string) (os.FileMode, error) {
	m, err := strconv.ParseUint(s, 8, 32)

	return os.FileMode(m), err
}

func createDirectory(ctx context.Context, p string) error {
	p, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	}

	return os.Mkdir(p, 0o700)
}

func removeFile(ctx context.Context, ty, p string, force bool) error {
	q, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	} else if _, err := os.Lstat(q); !force && errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s %q does not exist", ty, p)
	}

	return os.RemoveAll(q)
}

func runCommand(ctx context.Context, successfully, line, asynchronously string) (context.Context, error) {
	c, err := newCommand(line)
	if err != nil {
		return ctx, err
	}

	w := contextWorld(ctx)
	c.Cmd.Dir = w.CurrentDirectory
	c.Cmd.Env = w.Environment
	w = w.AddCommand(c)
	ctx = contextWithWorld(ctx, w)

	w.Stdin, err = c.Cmd.StdinPipe()
	if err != nil {
		return ctx, err
	}

	ctx = contextWithWorld(ctx, w)

	err = c.Cmd.Start()
	if err != nil {
		return ctx, err
	}

	time.Sleep(w.StartupWaitTime)

	if asynchronously == "" {
		if err := c.Cmd.Wait(); successfully != "" && err != nil {
			return ctx, fmt.Errorf("%v (stderr: %q)", err, w.Stderr())
		}
	} else if successfully != "" {
		return ctx, errors.New("cannot check an asynchronous command's success")
	}

	return ctx, nil
}

func runScript(ctx context.Context, s *godog.DocString) (context.Context, error) {
	const scriptPath = "script"

	err := createFile(ctx, scriptPath, s.Content)
	if err != nil {
		return ctx, err
	}

	return runCommand(ctx, "", "sh "+scriptPath, "")
}

func waitForStartup(ctx context.Context, seconds float64) context.Context {
	w := contextWorld(ctx)
	w.StartupWaitTime = time.Duration(seconds * float64(time.Second))

	return contextWithWorld(ctx, w)
}

func exitStatus(ctx context.Context, not string, code int) error {
	w := contextWorld(ctx)
	c := w.LastCommand()
	_ = c.Cmd.Wait()

	if c := c.Cmd.ProcessState.ExitCode(); (c == code) != (not == "") {
		return fmt.Errorf("expected exit code %d%s to be %d (stderr: %q)", c, not, code, w.Stderr())
	}

	return nil
}

func stdin(ctx context.Context, p string) error {
	w := contextWorld(ctx)

	q, err := w.path(p)
	if err != nil {
		return err
	}

	f, err := os.Open(q)
	if err != nil {
		return err
	}

	// TODO Figure out why we need to ignore errors...
	_, _ = io.Copy(w.Stdin, f)
	_ = w.Stdin.Close()

	return nil
}

func output(ctx context.Context, channel, from, not, exactly, pattern string) error {
	w := contextWorld(ctx)
	s := ""

	if from == "" {
		switch channel {
		case "stdout":
			s = w.Stdout()
		case "stderr":
			s = w.Stderr()
		default:
			s = w.Output()
		}
	} else {
		from := unescapeString(from)

		c := w.FindCommand(from)
		if c == nil {
			return fmt.Errorf("no command matching %q", from)
		}

		_ = c.Cmd.Wait()

		switch channel {
		case "stdout":
			s = c.Stdout.String()
		case "stderr":
			s = c.Stderr.String()
		default:
			s = c.Output()
		}
	}

	ok := strings.Contains(s, pattern)

	// Output from a specific command keeps its trailing newlines on exact matches like Aruba's.
	if exactly != "" && from != "" {
		ok = s == pattern
	} else if exactly != "" {
		ok = matchesExactly(s, pattern)
	}

	if ok != (not == "") {
		return fmt.Errorf("expected %s %q%s to contain%s %q", channel, s, not, exactly, pattern)
	}

	return nil
}

func fileContains(ctx context.Context, p, not, exactly, pattern string) error {
	q, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	}

	bs, err := os.ReadFile(q)
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
	q, err := contextWorld(ctx).path(p)
	if err != nil {
		return err
	}

	if i, err := os.Stat(q); (err == nil && i.IsDir() == (ty == "directory")) != (not == "") {
		return fmt.Errorf("%s %q should%s exist", ty, p, not)
	}

	return nil
}

func setEnvVar(ctx context.Context, k, v string) context.Context {
	w := contextWorld(ctx)
	w.Environment = append(w.Environment, k+"="+v)
	return contextWithWorld(ctx, w)
}

func appendEnvVar(ctx context.Context, v, k string) context.Context {
	return setEnvVar(ctx, k, contextWorld(ctx).EnvironmentVariable(k)+v)
}

func changeDirectory(ctx context.Context, p string) (context.Context, error) {
	w := contextWorld(ctx)

	d, err := w.path(p)
	if err != nil {
		return ctx, err
	}

	w.CurrentDirectory = d

	return contextWithWorld(ctx, w), nil
}

// [InitializeScenario] initializes a scenario.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(before)
	ctx.After(after)

	ctx.Step(
		`^a file named "(.*)" with "(.*)"$`,
		func(ctx context.Context, p, s string) error {
			return createFile(ctx, unescapeString(p), unescapeString(s))
		})
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.*)"(?: with mode "(.*)" and)? with:$`,
		func(ctx context.Context, p, mode string, s *godog.DocString) error {
			if mode == "" {
				return createFile(ctx, unescapeString(p), s.Content)
			}

			m, err := parseFileMode(unescapeString(mode))
			if err != nil {
				return err
			}

			return createFileWithMode(ctx, unescapeString(p), s.Content, m)
		})
	ctx.Step(
		`^(?:an|the) executable(?: named)? "(.*)" with:$`,
		func(ctx context.Context, p string, s *godog.DocString) error {
			return createFileWithMode(ctx, unescapeString(p), s.Content, 0o755)
		})
	ctx.Step(
		`^a directory named "(.*)"$`,
		func(ctx context.Context, p string) error {
			return createDirectory(ctx, unescapeString(p))
		})
	ctx.Step(
		`^(?:a|the) (directory|file)(?: named)? "(.*)" does not exist$`,
		func(ctx context.Context, ty, p string) error {
			return removeFile(ctx, ty, unescapeString(p), true)
		})
	ctx.Step(
		`^I remove (?:a|the) (directory|file)(?: named)? "(.*)"( with full force)?$`,
		func(ctx context.Context, ty, p, force string) error {
			return removeFile(ctx, ty, unescapeString(p), force != "")
		})
	ctx.Step("^I( successfully)? run `(.*)`( interactively| in the background)?$", runCommand)
	ctx.Step(`^I wait ([\d.]+) seconds? for (?:a|the) command to start up$`, waitForStartup)
	ctx.Step(`^the exit status should( not)? be (\d+)$`, exitStatus)
	ctx.Step(
		`^(?:the )?(output|std(?:out|err))(?: from "(.*)")? should( not)? contain( exactly)? "(.*)"$`,
		func(ctx context.Context, channel, from, not, exactly, pattern string) error {
			return output(ctx, channel, from, not, exactly, unescapeString(pattern))
		},
	)
	ctx.Step(
		`^(?:the )?(output|std(?:out|err))(?: from "(.*)")? should( not)? contain( exactly)?:$`,
		func(ctx context.Context, channel, from, not, exactly string, docString *godog.DocString) error {
			return output(ctx, channel, from, not, exactly, docString.Content)
		},
	)
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.*)" should( not)? contain "(.*)"$`,
		func(ctx context.Context, p, not, pattern string) error {
			return fileContains(ctx, unescapeString(p), not, "", unescapeString(pattern))
		})
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.*)" should( not)? contain( exactly)?:$`,
		func(ctx context.Context, p, not, exactly string, docString *godog.DocString) error {
			return fileContains(ctx, unescapeString(p), not, exactly, trimTrailingNewlines(docString.Content))
		})
	ctx.Step(
		`^I pipe in the file(?: named)? "(.*)"$`,
		func(ctx context.Context, p string) error {
			return stdin(ctx, unescapeString(p))
		})
	ctx.Step(
		`^(?:a|the) (directory|file)(?: named)? "(.*)" should( not)? exist$`,
		func(ctx context.Context, ty, p, not string) error {
			return fileExists(ctx, ty, unescapeString(p), not)
		})
	ctx.Step(
		`^I set the environment variable "(.*)" to "(.*)"$`,
		func(ctx context.Context, k, v string) context.Context {
			return setEnvVar(ctx, unescapeString(k), unescapeString(v))
		})
	ctx.Step(
		`^I append "(.*)" to the environment variable "(.*)"$`,
		func(ctx context.Context, v, k string) context.Context {
			return appendEnvVar(ctx, unescapeString(v), unescapeString(k))
		})
	ctx.Step(`^I run the following (?:commands|script):$`, runScript)
	ctx.Step(
		`^I cd to "(.*)"$`,
		func(ctx context.Context, p string) (context.Context, error) {
			return changeDirectory(ctx, unescapeString(p))
		})
}
