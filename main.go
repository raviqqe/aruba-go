package main

import (
	"os"

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

func main() {
	pflag.Parse()
	options.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                 "godog",
		TestSuiteInitializer: func(*godog.TestSuiteContext) {},
		ScenarioInitializer:  func(*godog.ScenarioContext) {},
		Options:              &options,
	}.Run()

	os.Exit(status)
}
