package main_test

import (
	"testing"

	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

func TestRunNoFeature(t *testing.T) {
	status, err := main.Run()

	assert.Equal(t, 1, status)
	assert.Equal(t, nil, err)
}
