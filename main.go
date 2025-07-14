package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

var options = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
}

func init() {
	godog.BindCommandLineFlags("", &options)
}

func InitializeScenario(scenario *godog.ScenarioContext) {
	scenario.Step(`^a file named {string} with:$`, func() {})

	scenario.Step("^I (successfully |)run `(.*)`$", func(ctx context.Context, successfully, command string) (context.Context, error) {
		ss := strings.Split(command, " ")
		c := exec.Command(ss[0], ss[1:]...)
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
