package aruba_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raviqqe/aruba-go"
	"github.com/stretchr/testify/assert"
)

func TestSuccessFeatures(t *testing.T) {
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
