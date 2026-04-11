package main

import (
	"testing"

	"github.com/AbsaOSS/golic/cmd/commands"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdErrorMessage(t *testing.T) {
	cmd, _ := commands.SetupTest()

	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.ErrorContains(t, err, "no arguments passed")
}
