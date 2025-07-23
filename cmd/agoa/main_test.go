package main_test

import (
	"testing"

	"github.com/raviqqe/aruba-go/cmd/agoa"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	status, err := main.Run()

	assert.Equal(t, 1, status)
	assert.Equal(t, nil, err)
}
