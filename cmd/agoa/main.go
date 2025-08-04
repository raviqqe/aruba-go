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

const Version = "0.1.7"

type Options struct {
	Godog   godog.Options
	Version bool
}

func main() {
	status, err := Run(parseOptions())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(status)

}

func Run(options Options) (int, error) {
	if options.Version {
		_, err := fmt.Fprintln(options.Godog.Output, Version)

		if err != nil {
			return 1, err
		}

		return 0, nil
	}

	suite := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options:             &options.Godog,
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

func parseOptions() Options {
	options := Options{
		Godog: godog.Options{
			Concurrency: runtime.NumCPU(),
			Format:      "pretty",
			Output:      colors.Colored(os.Stdout),
			Strict:      true,
		},
	}

	godog.BindCommandLineFlags("", &options.Godog)
	pflag.BoolVar(&options.Version, "version", false, "show version")

	pflag.Parse()
	options.Godog.Paths = pflag.Args()

	return options
}
