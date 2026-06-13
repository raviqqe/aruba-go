package aruba

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/kballard/go-shellquote"
)

func parseString(s string) (string, error) {
	return strconv.Unquote(s)
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
	err := os.RemoveAll(contextWorld(ctx).RootDirectory)

	return ctx, err
}

func worldPath(w world, p string) (string, error) {
	q := filepath.Join(w.CurrentDirectory, p)

	if d, err := filepath.Rel(w.RootDirectory, q); err != nil {
		return "", err
	} else if strings.HasPrefix(d, "..") {
		return "", fmt.Errorf("path %q is outside the working directory", p)
	}

	return q, nil
}

func createFile(ctx context.Context, p, s string) error {
	p, err := worldPath(contextWorld(ctx), p)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(p), 0o700); err != nil {
		return err
	}

	return os.WriteFile(p, []byte(s), 0o600)
}

func createFileWithMode(ctx context.Context, p, s string, mode os.FileMode) error {
	if err := createFile(ctx, p, s); err != nil {
		return err
	}

	q, err := worldPath(contextWorld(ctx), p)
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
	p, err := worldPath(contextWorld(ctx), p)
	if err != nil {
		return err
	}

	return os.Mkdir(p, 0o700)
}

func runCommand(ctx context.Context, successfully, command, interactively string) (context.Context, error) {
	command, err := parseString(command)
	if err != nil {
		return ctx, err
	}

	ss, err := shellquote.Split(command)
	if err != nil {
		return ctx, err
	} else if len(ss) == 0 {
		return ctx, errors.New("empty command")
	}

	c := exec.Command(ss[0], ss[1:]...)
	w := contextWorld(ctx)
	c.Dir = w.CurrentDirectory
	c.Stdout = bytes.NewBuffer(nil)
	c.Stderr = bytes.NewBuffer(nil)
	c.Env = w.Environment
	w = w.AddCommand(c)
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
		if err := c.Wait(); successfully != "" && err != nil {
			return ctx, fmt.Errorf("%v (stderr: %q)", err, w.Stderr())
		}
	} else if successfully != "" {
		return ctx, errors.New("cannot check interactive command's success")
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
	w := contextWorld(ctx)
	c := w.LastCommand()
	_ = c.Wait()

	if c := c.ProcessState.ExitCode(); (c == code) != (not == "") {
		return fmt.Errorf("expected exit code %d%s to be %d (stderr: %q)", c, not, code, w.Stderr())
	}

	return nil
}

func stdin(ctx context.Context, p string) error {
	w := contextWorld(ctx)

	q, err := worldPath(w, p)
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
		from, err := parseString(from)
		if err != nil {
			return err
		}

		c := w.FindCommand(from)
		_ = c.Wait()

		switch channel {
		case "stdout":
			s = c.Stdout.(*bytes.Buffer).String()
		case "stderr":
			s = c.Stderr.(*bytes.Buffer).String()
		default:
			s = c.Stdout.(*bytes.Buffer).String() + c.Stderr.(*bytes.Buffer).String()
		}
	}

	if exactly == "" && strings.Contains(s, pattern) != (not == "") ||
		exactly != "" && matchesExactly(s, pattern) != (not == "") {
		return fmt.Errorf("expected %s %q%s to contain%s %q", channel, s, not, exactly, pattern)
	}

	return nil
}

func fileContains(ctx context.Context, p, not, exactly, pattern string) error {
	q, err := worldPath(contextWorld(ctx), p)
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
	q, err := worldPath(contextWorld(ctx), p)
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

func changeDirectory(ctx context.Context, p string) (context.Context, error) {
	w := contextWorld(ctx)

	d, err := worldPath(w, p)
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
		`^a file named "(.+)" with (".*")$`,
		func(ctx context.Context, p, s string) error {
			s, err := parseString(s)
			if err != nil {
				return err
			}
			return createFile(ctx, p, s)
		})
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.+)"(?: with mode "(.*)" and)? with:$`,
		func(ctx context.Context, p, mode string, s *godog.DocString) error {
			content := trimTrailingNewlines(s.Content) + "\n"

			if mode == "" {
				return createFile(ctx, p, content)
			}

			m, err := parseFileMode(mode)
			if err != nil {
				return err
			}

			return createFileWithMode(ctx, p, content, m)
		})
	ctx.Step(
		`^(?:an|the) executable(?: named)? "(.+)" with:$`,
		func(ctx context.Context, p string, s *godog.DocString) error {
			return createFileWithMode(ctx, p, trimTrailingNewlines(s.Content)+"\n", 0o755)
		})
	ctx.Step(`^a directory named "(.+)"$`, createDirectory)
	ctx.Step("^I( successfully)? run (`.*`)( interactively)?$", runCommand)
	ctx.Step(`^the exit status should( not)? be (\d+)$`, exitStatus)
	ctx.Step(
		`^the (output|std(?:out|err))(?: from (".*"))? should( not)? contain( exactly)? (".*")$`,
		func(ctx context.Context, channel, from, not, exactly, pattern string) error {
			pattern, err := parseString(pattern)
			if err != nil {
				return err
			}

			return output(ctx, channel, from, not, exactly, pattern)
		},
	)
	ctx.Step(
		`^the (output|std(?:out|err))(?: from (".*"))? should( not)? contain( exactly)?:$`,
		func(ctx context.Context, channel, from, not, exactly string, docString *godog.DocString) error {
			return output(ctx, channel, from, not, exactly, trimTrailingNewlines(docString.Content))
		},
	)
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.*)" should( not)? contain (".*")$`,
		func(ctx context.Context, p, not, pattern string) error {
			pattern, err := parseString(pattern)
			if err != nil {
				return err
			}

			return fileContains(ctx, p, not, "", pattern)
		})
	ctx.Step(
		`^(?:a|the) file(?: named)? "(.+)" should( not)? contain( exactly)?:$`,
		func(ctx context.Context, p, not, exactly string, docString *godog.DocString) error {
			return fileContains(ctx, p, not, exactly, trimTrailingNewlines(docString.Content))
		})
	ctx.Step(`^I pipe in the file(?: named)? "(.*)"$`, stdin)
	ctx.Step(`^(?:a|the) (directory|file)(?: named)? "(.*)" should( not)? exist$`, fileExists)
	ctx.Step(`^I set the environment variable "(.*)" to "(.*)"$`, setEnvVar)
	ctx.Step(`^I run the following (?:commands|script):$`, runScript)
	ctx.Step(`^I cd to "(.*)"$`, changeDirectory)
}
