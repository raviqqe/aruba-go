package aruba

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStopKillsRunningCommand(t *testing.T) {
	c := exec.Command("sleep", "10")
	require.NoError(t, c.Start())

	world{}.AddCommand(c).Stop()

	require.NotNil(t, c.ProcessState)
	assert.False(t, c.ProcessState.Success())
}

func TestStopKeepsFinishedCommand(t *testing.T) {
	c := exec.Command("true")
	require.NoError(t, c.Start())
	require.NoError(t, c.Wait())

	world{}.AddCommand(c).Stop()

	assert.True(t, c.ProcessState.Success())
}

func TestStopIgnoresUnstartedCommand(t *testing.T) {
	c := exec.Command("true")

	world{}.AddCommand(c).Stop()

	assert.Nil(t, c.ProcessState)
}
