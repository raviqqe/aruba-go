package main_test

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

var defaultOptions = main.Options{
	Godog: godog.Options{
		Concurrency: runtime.NumCPU(),
		Format:      "pretty",
		Output:      colors.Colored(os.Stdout),
		Strict:      true,
	},
}

func TestRunNoFeature(t *testing.T) {
	b := bytes.NewBuffer(nil)
	options := defaultOptions
	options.Godog.Output = b
	status, err := main.Run(options)

	assert.Equal(t, 1, status)
	assert.Nil(t, err)

	s := b.String()

	assert.Regexp(t, "No scenarios", s)
	assert.Regexp(t, "No steps", s)
}

func TestRunFeatures(t *testing.T) {
	err := os.Chdir("../..")
	assert.Nil(t, err)

	options := defaultOptions
	options.Godog.Output = io.Discard
	status, err := main.Run(options)

	assert.Zero(t, status)
	assert.Nil(t, err)
}

func TestRunVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)
	options := defaultOptions
	options.Godog.Output = b
	options.Version = true
	status, err := main.Run(options)

	assert.Equal(t, 0, status)
	assert.Nil(t, err)

	assert.Equal(t, main.Version+"\n", b.String())
}
