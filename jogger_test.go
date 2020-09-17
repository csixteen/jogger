package jogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandSuccess(t *testing.T) {
	_, _, err := Run("ls", []string{"-la"})
	assert.Nil(t, err)
}

func TestRunCommandNoOutput(t *testing.T) {
	_, _, err := Run("ls", []string{"-la"}, NoOutput())
	assert.Nil(t, err)
}

func TestRunCommandFailure(t *testing.T) {
	_, _, err := Run("lssssss", []string{"-la"})
	assert.Error(t, err)
}
