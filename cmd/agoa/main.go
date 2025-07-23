package main

import (
	"errors"
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
	pflag.Parse()
	options.Paths = pflag.Args()

	suite := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options:             &options,
	}

	fs, err := suite.RetrieveFeatures()
	if err != nil {
		fail(err)
	} else if len(fs) == 0 {
		fail(errors.New("no features found"))
	}

	os.Exit(suite.Run())
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
