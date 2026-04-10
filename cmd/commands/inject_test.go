package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInject(t *testing.T) {
	root, out := SetupTest()

	t.Run("inject with missing config file", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject"})
		err := root.Execute()
		assert.ErrorContains(t, err, "ensure '.golic.yaml' exists")
	})

	t.Run("inject with missing ignore file", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject", "-p", "../../.golic.yaml"})
		err := root.Execute()
		assert.ErrorContains(t, err, "ensure '.licignore' exists")
	})

	t.Run("inject with missing template", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject", "-p", "../../.golic.yaml", "-l", "../../.licignore"})
		err := root.Execute()
		assert.ErrorContains(t, err, "licence template not provided")
	})

	t.Run("inject with mit", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject", "-p", "../../.golic.yaml", "-l", "../../.licignore", "-t", "mit"})
		_ = root.Execute()
	})

}
