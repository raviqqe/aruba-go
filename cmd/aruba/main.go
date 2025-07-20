package main

import (
	"os"
	"runtime"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raviqqe/aruba-go"
	"github.com/spf13/pflag"
)

type directoryKey struct{}
type exitCodeKey struct{}
type stdoutKey struct{}
type stderrKey struct{}

var options = godog.Options{
	Concurrency: runtime.NumCPU(),
	Output:      colors.Colored(os.Stdout),
	Format:      "pretty",
}

func init() {
	godog.BindCommandLineFlags("", &options)
}

func main() {
	pflag.Parse()
	options.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options:             &options,
	}.Run()

	os.Exit(status)
}
