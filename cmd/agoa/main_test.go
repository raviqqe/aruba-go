package main_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

var defaultOptions = main.Options{}

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
