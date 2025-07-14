package main

import (
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

func InitializeScenario(context *godog.ScenarioContext) {
	context.Step(`^a file named {string} with:$`, func() {})
	context.Step("^I (successfully |)run `(.*)`$", func(successfully, command string) error {
		components := strings.Split(" ", command)

		return exec.Command(components[0], components[1:]...).Run()
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
