package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmdErrorMessage(t *testing.T) {
	cmd := RootCmd()
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs([]string{})

	_ = cmd.Execute()

	assert.Contains(t, b.String(), "Error: no parameters included")
}
