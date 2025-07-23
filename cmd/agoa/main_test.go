package main_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

func TestRunNoFeature(t *testing.T) {
	b := bytes.NewBuffer(nil)
	status, err := main.Run(b, true)

	assert.Equal(t, 1, status)
	assert.Equal(t, nil, err)

	s := b.String()

	assert.Regexp(t, `No scenarios`, s)
	assert.Regexp(t, `No steps`, s)
}

func TestRunFeatures(t *testing.T) {
	os.Chdir("../..")
	status, err := main.Run(io.Discard, true)

	assert.Equal(t, 0, status)
	assert.Equal(t, nil, err)
}
