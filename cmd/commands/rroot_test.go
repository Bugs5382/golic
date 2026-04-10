package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	root, out := SetupTest()

	t.Run("root no args passed", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{})
		err := root.Execute()
		assert.ErrorContains(t, err, "no arguments passed")
	})

}
