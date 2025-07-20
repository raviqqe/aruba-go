package aruba

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
)

type directoryKey struct{}
type exitCodeKey struct{}
type stdoutKey struct{}
type stderrKey struct{}

var doubleQuotePattern = regexp.MustCompile(`([^\\])"`)
var backSlashPattern = regexp.MustCompile(`([^\\])\\`)
var headDoubleQuotePattern = regexp.MustCompile(`^"`)

func quote(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = doubleQuotePattern.ReplaceAllString(s, `$1\"`)
	s = headDoubleQuotePattern.ReplaceAllString(s, `\"`)

	return s
}

var unquotePattern = regexp.MustCompile(`\\(\\)`)

func unquote(s string) string {
	return unquotePattern.ReplaceAllString(s, `$1`)
}

func parseString(s string) string {
	return unquote(strings.TrimSpace(s))
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "godog-*")

	return context.WithValue(ctx, directoryKey{}, d), err
}

func createFile(ctx context.Context, p string, docString *godog.DocString) error {
	return os.WriteFile(
		path.Join(ctx.Value(directoryKey{}).(string), p),
		[]byte(unquote(docString.Content)),
		0o644,
	)
}

func runCommand(ctx context.Context, successfully, command string) (context.Context, error) {
	// TODO Unquote only once?
	command = unquote(parseString(command))

	ss := strings.Split(command, " ")
	c := exec.Command(ss[0], ss[1:]...)
	c.Dir = ctx.Value(directoryKey{}).(string)
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout
	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr

	err := c.Run()
	ctx = context.WithValue(ctx, exitCodeKey{}, c.ProcessState.ExitCode())
	ctx = context.WithValue(ctx, stdoutKey{}, stdout.Bytes())
	ctx = context.WithValue(ctx, stderrKey{}, stderr.Bytes())

	if successfully == "" {
		return ctx, nil
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

	if exactly == "" && strings.Contains(quote(s), pattern) != (not == "") {
		return fmt.Errorf("expected %s%s to contain %q but got %q", stdout, not, pattern, s)
	} else if exactly != "" && (quote(s) == pattern || quote(strings.TrimSpace(s)) == pattern) != (not == "") {
		return fmt.Errorf("expected %s%s to be %q but got %q", stdout, not, pattern, s)
	}

	return nil
}

func fileContains(ctx context.Context, p, not, exactly, pattern string) error {
	bs, err := os.ReadFile(path.Join(ctx.Value(directoryKey{}).(string), p))
	if err != nil {
		return err
	}

	s := string(bs)
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
			return stdout(ctx, port, not, exactly, parseString(pattern))
		},
	)
	ctx.Step(
		`^the (std(?:out|err)) should( not)? contain( exactly)?:$`,
		func(ctx context.Context, port, not, exactly string, docString *godog.DocString) error {
			return stdout(ctx, port, not, exactly, quote(strings.TrimSpace(docString.Content)))
		},
	)
	ctx.Step(`^a file named "([^"]*)" should( not)? contain "([^"]*)"$`, func(ctx context.Context, p, not, pattern string) error {
		return fileContains(ctx, p, not, "", parseString(pattern))
	})
	ctx.Step(`^a file named "([^"]*)" should( not)? contain( exactly)?:$`, func(ctx context.Context, p, not, exactly string, docString *godog.DocString) error {
		return fileContains(ctx, p, not, exactly, strings.TrimSpace(docString.Content))
	})
}
