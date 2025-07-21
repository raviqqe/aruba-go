package aruba

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type directoryKey struct{}
type exitCodeKey struct{}
type stdoutKey struct{}
type stderrKey struct{}

func unquote(s string) (string, error) {
	s, err := strconv.Unquote(`"` + s + `"`)

	if err != nil {
		return "", err
	}

	return s, nil
}

func parseString(s string) (string, error) {
	return unquote(strings.TrimSpace(s))
}

func parseDocString(s string) string {
	return strings.TrimSpace(s)
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "aruba-*")

	return context.WithValue(ctx, directoryKey{}, d), err
}

func createFile(ctx context.Context, p string, docString *godog.DocString) error {
	return os.WriteFile(
		path.Join(ctx.Value(directoryKey{}).(string), p),
		[]byte(parseDocString(docString.Content)),
		0o644,
	)
}

func runCommand(ctx context.Context, successfully, command string) (context.Context, error) {
	command, err := parseString(command)
	if err != nil {
		return ctx, err
	}

	ss := strings.Split(command, " ")
	c := exec.Command(ss[0], ss[1:]...)
	c.Dir = ctx.Value(directoryKey{}).(string)
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout
	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr

	err = c.Run()
	ctx = context.WithValue(ctx, exitCodeKey{}, c.ProcessState.ExitCode())
	ctx = context.WithValue(ctx, stdoutKey{}, stdout.Bytes())
	ctx = context.WithValue(ctx, stderrKey{}, stderr.Bytes())

	if successfully == "" {
		err = nil
	}

	return ctx, err
}

func exitStatus(ctx context.Context, not string, code int) error {
	if c := ctx.Value(exitCodeKey{}).(int); (c == code) != (not == "") {
		return fmt.Errorf("expected exit code%s %d but got %d", not, code, c)
	}

	return nil
}

func stdout(ctx context.Context, stdout, not, exactly, pattern string) error {
	key := any(stdoutKey{})

	if stdout == "stderr" {
		key = stderrKey{}
	}

	s := string(ctx.Value(key).([]byte))

	if exactly == "" && strings.Contains(s, pattern) != (not == "") ||
		exactly != "" && (s == pattern || strings.TrimSpace(s) == pattern) != (not == "") {
		return fmt.Errorf("expected %s %q%s to contain%s %q", stdout, s, not, exactly, pattern)
	}

	return nil
}

func fileContains(ctx context.Context, p, not, exactly, pattern string) error {
	bs, err := os.ReadFile(path.Join(ctx.Value(directoryKey{}).(string), p))
	if err != nil {
		return err
	}

	s := strings.TrimSpace(string(bs))
	ok := strings.Contains(s, pattern)

	if exactly != "" {
		ok = s == pattern
	}

	if ok != (not == "") {
		return fmt.Errorf("expected file %q%s to contain %q", p, not, pattern)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(before)
	ctx.Step(`^a file named "((?:\\.|[^"\\])+)" with:$`, createFile)
	ctx.Step("^I( successfully)? run `(.*)`$", runCommand)
	ctx.Step(`^the exit status should( not)? be (\d+)$`, exitStatus)
	ctx.Step(
		`^the (std(?:out|err)) should( not)? contain( exactly)? "((?:\\.|[^"\\])*)"$`,
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
	ctx.Step(`^a file named "([^"]*)" should( not)? contain "([^"]*)"$`, func(ctx context.Context, p, not, pattern string) error {
		pattern, err := parseString(pattern)
		if err != nil {
			return err
		}

		return fileContains(ctx, p, not, "", pattern)
	})
	ctx.Step(`^a file named "([^"]*)" should( not)? contain( exactly)?:$`, func(ctx context.Context, p, not, exactly string, docString *godog.DocString) error {
		return fileContains(ctx, p, not, exactly, parseDocString(docString.Content))
	})
}
