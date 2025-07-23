package aruba_test

import (
	"bytes"
	"os"
	"regexp"
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
			Concurrency: 1,
			Format:      "pretty",
			NoColors:    true,
			Output:      b,
			Paths:       []string{"failures"},
		},
	}.Run()

	assert.NotZero(t, status)

	snaps.MatchSnapshot(
		t,
		regexp.MustCompile(` *# scenario\.go.*`).
			ReplaceAllString(
				regexp.MustCompile(`[[:space:]]*[0-9.]+ms`).
					ReplaceAllString(b.String(), ""),
				""),
	)
}
