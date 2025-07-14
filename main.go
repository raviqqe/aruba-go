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

func InitializeScenario(scenario *godog.ScenarioContext) {
	scenario.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		d, err := os.MkdirTemp("", "godog-*")

		return context.WithValue(ctx, "directory", d), err
	})

	scenario.Step(`^a file named "((\\\\|\\"|[^"\\])+)" with:$`, func(ctx context.Context, p string) error {
		return os.WriteFile(path.Join(ctx.Value("directory").(string), p), nil, 0o644)
	})

	scenario.Step("^I (successfully |)run `(.*)`$", func(ctx context.Context, successfully, command string) (context.Context, error) {
		ss := strings.Split(command, " ")
		c := exec.Command(ss[0], ss[1:]...)
		c.Dir = ctx.Value("directory").(string)

		err := c.Run()
		ctx = context.WithValue(ctx, "exitCode", c.ProcessState.ExitCode())

		if successfully == "" {
			return ctx, nil
		}

		return ctx, err
	})

	scenario.Step(`^the exit status should (not |)be (\d+)$`, func(ctx context.Context, not string, code int) error {
		c := ctx.Value("exitCode").(int)

		if not == "" && c == code {
			return nil
		}

		return fmt.Errorf("expected exit code %d but got %d", code, c)
	})
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
