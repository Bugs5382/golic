package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	root, out := SetupTest()

	t.Run("remove with missing config file", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"remove"})
		err := root.Execute()
		assert.ErrorContains(t, err, "ensure '.golic.yaml' exists")
	})

	t.Run("remove with missing ignore file", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"remove", "-p", "../../.golic.yaml"})
		err := root.Execute()
		assert.ErrorContains(t, err, "ensure '.licignore' exists")
	})

	t.Run("remove with missing template", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"remove", "-p", "../../.golic.yaml", "-l", "../../.licignore"})
		err := root.Execute()
		assert.ErrorContains(t, err, "licence template not provided")
	})

}
