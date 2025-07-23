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

var options = godog.Options{
	Concurrency: runtime.NumCPU(),
	Output:      colors.Colored(os.Stdout),
	Format:      "pretty",
	Strict:      true,
}

func init() {
	godog.BindCommandLineFlags("", &options)
}

func main() {
	status, err := Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(status)

}

func Run() (int, error) {
	pflag.Parse()
	options.Paths = pflag.Args()

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
