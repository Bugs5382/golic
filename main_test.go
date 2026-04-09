package main

import (
	"testing"

	"github.com/AbsaOSS/golic/cmd/commands"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdErrorMessage(t *testing.T) {
	cmd, b := commands.SetupTest()

	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	assert.Contains(t, b.String(), "Error: no parameters included")
}
