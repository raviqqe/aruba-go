package main

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
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

type directoryKey struct{}
type exitCodeKey struct{}
type stdoutKey struct{}
type stderrKey struct{}

var options = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

func init() {
	godog.BindCommandLineFlags("", &options)
}

func quote(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\"", "\\\"")

	return s
}

var normalUnquotePattern = regexp.MustCompile(`\\(\\|"|n)`)

func unquote(s string) string {
	return normalUnquotePattern.ReplaceAllStringFunc(s, func(s string) string {
		switch s {
		case "\\n":
			return "\n"
		default:
			return s[1:]
		}
	})
}

var simpleUnquotePattern = regexp.MustCompile(`\\(\\)`)

func unquoteSimple(s string) string {
	return simpleUnquotePattern.ReplaceAllString(s, `$1`)
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "godog-*")

	return context.WithValue(ctx, directoryKey{}, d), err
}

func createFile(ctx context.Context, p string, docString *godog.DocString) error {
	return os.WriteFile(
		path.Join(ctx.Value(directoryKey{}).(string), p),
		[]byte(unquoteSimple(docString.Content)),
		0o644,
	)
}

func runCommand(ctx context.Context, successfully, command string) (context.Context, error) {
	// TODO Unquote only once?
	command = unquote(unquote(command))

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
	c := ctx.Value(exitCodeKey{}).(int)

	if (not == "") == (c == code) {
		return nil
	}

	return fmt.Errorf("expected exit code %s%d but got %d", not, code, c)
}

func stdout(ctx context.Context, stdout, not, exactly, expected string) error {
	key := any(stdoutKey{})

	if stdout == "stderr" {
		key = stderrKey{}
	}

	s := string(ctx.Value(key).([]byte))
	expected = unquoteSimple(strings.TrimSpace(expected))

	if exactly == "" && !strings.Contains(quote(s), expected) {
		return fmt.Errorf("expected %s to contain %q but got %q", stdout, expected, s)
	} else if exactly != "" {
		if t := strings.TrimSpace(s); t != "" {
			s = t
		}
		s := quote(s)

		if s != expected {
			return fmt.Errorf("expected %s to be %q but got %q", stdout, expected, s)
		}
	}

	return nil
}

func initializeScenario(scenario *godog.ScenarioContext) {
	scenario.Before(before)
	scenario.Step(`^a file named "((?:\\.|[^"\\])+)" with:$`, createFile)
	scenario.Step("^I (successfully |)run `(.*)`$", runCommand)
	scenario.Step(`^the exit status should (not |)be (\d+)$`, exitStatus)
	scenario.Step(`^the (std(?:out|err)) should (not |)contain (exactly |)"((?:\\.|[^"\\])+)"$`, stdout)
}

func main() {
	pflag.Parse()
	options.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "godog",
		ScenarioInitializer: initializeScenario,
		Options:             &options,
	}.Run()

	os.Exit(status)
}
