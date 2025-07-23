package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raviqqe/aruba-go"
	"github.com/spf13/pflag"
)

func main() {
	status, err := Run(os.Stdout, false)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(status)

}

func Run(out io.Writer, test bool) (int, error) {
	options := parseOptions(out, test)

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

func parseOptions(out io.Writer, test bool) godog.Options {
	options := godog.Options{
		Concurrency: runtime.NumCPU(),
		Output:      colors.Colored(out),
		Format:      "pretty",
		Strict:      true,
	}

	if !test {
		godog.BindCommandLineFlags("", &options)
	}

	pflag.Parse()
	options.Paths = pflag.Args()

	return options
}
