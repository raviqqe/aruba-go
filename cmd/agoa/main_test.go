package main_test

import (
	"os"
	"testing"

	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

func TestRunNoFeature(t *testing.T) {
	status, err := main.Run()

	assert.Equal(t, 1, status)
	assert.Equal(t, nil, err)
}

func TestRunFeatures(t *testing.T) {
	os.Chdir("../..")
	status, err := main.Run()

	assert.Equal(t, 0, status)
	assert.Equal(t, nil, err)
}
