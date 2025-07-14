package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
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

func before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	d, err := os.MkdirTemp("", "godog-*")

	return context.WithValue(ctx, "directory", d), err
}

func createFile(ctx context.Context, p string) error {
	return os.WriteFile(path.Join(ctx.Value("directory").(string), p), nil, 0o644)
}

func runCommand(ctx context.Context, successfully, command string) (context.Context, error) {
	ss := strings.Split(command, " ")
	c := exec.Command(ss[0], ss[1:]...)
	c.Dir = ctx.Value("directory").(string)

	err := c.Run()
	ctx = context.WithValue(ctx, "exitCode", c.ProcessState.ExitCode())

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

func InitializeScenario(scenario *godog.ScenarioContext) {
	scenario.Before(before)
	scenario.Step(`^a file named "((\\\\|\\"|[^"\\])+)" with:$`, createFile)
	scenario.Step("^I (successfully |)run `(.*)`$", runCommand)
	scenario.Step(`^the exit status should (not |)be (\d+)$`, exitStatus)
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
