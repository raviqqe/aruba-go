package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raviqqe/aruba-go"
	"github.com/spf13/pflag"
)

type options struct {
	godog   godog.Options
	version bool
}

func main() {
	status, err := Run(parseOptions())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(status)

}

func Run(options options) (int, error) {
	suite := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options:             &options.godog,
	}

	fs, err := suite.RetrieveFeatures()
	if err != nil {
		return 1, err
	}

	status := suite.Run()

	if len(fs) == 0 {
		status = 1
	}

	return status, nil
}

func parseOptions() options {
	options := options{
		godog: godog.Options{
			Concurrency: runtime.NumCPU(),
			Format:      "pretty",
			Output:      colors.Colored(os.Stdout),
			Strict:      true,
		},
	}

	godog.BindCommandLineFlags("", &options.godog)
	version := pflag.Bool("version", false, "Show version.")

	pflag.Parse()
	options.Paths = pflag.Args()

	return options
}
