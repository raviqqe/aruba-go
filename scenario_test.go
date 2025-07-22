package aruba_test

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/raviqqe/aruba-go"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulFeatures(t *testing.T) {
	status := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options: &godog.Options{
			Concurrency: runtime.NumCPU(),
			Output:      colors.Colored(os.Stdout),
			Format:      "pretty",
		},
	}.Run()

	assert.Zero(t, status)
}

func TestFailedFeatures(t *testing.T) {
	b := bytes.NewBuffer(nil)
	status := godog.TestSuite{
		Name:                "aruba",
		ScenarioInitializer: aruba.InitializeScenario,
		Options: &godog.Options{
			Concurrency: runtime.NumCPU(),
			Format:      "pretty",
			Output:      b,
			Paths:       []string{"failures"},
		},
	}.Run()

	assert.NotZero(t, status)
	snaps.MatchSnapshot(t, b.String())
}
