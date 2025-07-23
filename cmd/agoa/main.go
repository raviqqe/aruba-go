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

var DefaultOptions = godog.Options{
	Concurrency: runtime.NumCPU(),
	Format:      "pretty",
	Output:      colors.Colored(os.Stdout),
	Strict:      true,
}

func main() {
	status, err := Run(options())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(status)

}

func Run(options godog.Options) (int, error) {
	suite := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options:             &options,
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

func options() godog.Options {
	options := DefaultOptions

	godog.BindCommandLineFlags("", &options)

	pflag.Parse()
	options.Paths = pflag.Args()

	return options
}
