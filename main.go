package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

var options = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

func init() {
	godog.BindCommandLineFlags("", &options)
}

func unquote(s string) (string, error) {
	return strconv.Unquote(`"` + s + `"`)
}

var docStringPattern = regexp.MustCompile(`\\(\\|")`)

func unquoteDocString(s string) string {
	return docStringPattern.ReplaceAllString(s, `$1`)
}

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "godog-*")

	return context.WithValue(ctx, "directory", d), err
}

func createFile(ctx context.Context, p string, docString *godog.DocString) error {
	return os.WriteFile(
		path.Join(ctx.Value("directory").(string), p),
		[]byte(unquoteDocString(docString.Content)),
		0o644,
	)
}

func runCommand(ctx context.Context, successfully, command string) (context.Context, error) {
	command, err := unquote(command)
	if err != nil {
		return ctx, err
	}
	// TODO Unquote only once?
	command, err = unquote(command)
	if err != nil {
		return ctx, err
	}

	ss := strings.Split(command, " ")
	c := exec.Command(ss[0], ss[1:]...)
	c.Dir = ctx.Value("directory").(string)
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout
	stderr := bytes.NewBuffer(nil)
	c.Stderr = stderr

	err = c.Run()
	ctx = context.WithValue(ctx, "exitCode", c.ProcessState.ExitCode())
	ctx = context.WithValue(ctx, "stdout", stdout.Bytes())
	ctx = context.WithValue(ctx, "stderr", stderr.Bytes())

	if successfully == "" {
		return ctx, nil
	}

	return ctx, err
}

func exitStatus(ctx context.Context, not string, code int) error {
	c := ctx.Value("exitCode").(int)

	if (not == "") == (c == code) {
		return nil
	}

	return fmt.Errorf("expected exit code %s%d but got %d", not, code, c)
}

func stdout(ctx context.Context, stdout, not, exactly, expected string) error {
	s := string(ctx.Value(stdout).([]byte))
	expected = strings.TrimSpace(expected)
	expected, err := unquote(expected)
	if err != nil {
		return err
	}

	if exactly == "" && !strings.Contains(s, expected) {
		return fmt.Errorf("expected %s to contain %q but got %q", stdout, expected, s)
	} else if exactly != "" {
		s := strings.TrimSpace(s)

		if s != expected {
			return fmt.Errorf("expected %s to be %q but got %q", stdout, expected, s)
		}
	}

	return nil
}

func InitializeScenario(scenario *godog.ScenarioContext) {
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
		Name:                 "godog",
		TestSuiteInitializer: func(*godog.TestSuiteContext) {},
		ScenarioInitializer:  InitializeScenario,
		Options:              &options,
	}.Run()

	os.Exit(status)
}
