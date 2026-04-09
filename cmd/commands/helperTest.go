package commands

import (
	"bytes"

	"github.com/spf13/cobra"
)

func SetupTest() (*cobra.Command, *bytes.Buffer) {
	// 1. Get a fresh instance of the command
	root := RootCmd()
	b := new(bytes.Buffer)

	// 2. Redirect output and error streams to the buffer
	root.SetOut(b)
	root.SetErr(b)

	// 3. Reset arguments to ensure no leakage from other tests
	root.SetArgs([]string{})

	return root, b
}
