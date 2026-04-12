package commands

/*
Apache License 2.0

Copyright 2006 Shane

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInject(t *testing.T) {
	root, out := SetupTest()

	t.Run("inject with missing ignore file", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject"})
		err := root.Execute()
		assert.ErrorContains(t, err, "ensure '.licignore' exists")
	})

	t.Run("inject with missing template", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject", "-l", "../../.licignore"})
		err := root.Execute()
		assert.ErrorContains(t, err, "licence template not provided")
	})

	t.Run("inject with mit", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"inject", "-p", "../../.golic.yaml", "-l", "../../.licignore", "-t", "mit", "-d"})
		_ = root.Execute()
	})

}
