package commands

/*
Apache License 2.0

Copyright 2026 Shane

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
	"bytes"
	"os"
	"testing"

	"github.com/Bugs5382/golic"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {

	_ = os.Chdir(golic.GetProjectRoot())

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	t.Run("remove -- template missing", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"inject",
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "license template not provided")
	})

	t.Run("remove -- custom config file not found", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"remove",
			"-p",
			".golic-test.yaml",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "custom config file not found: .golic-test.yaml")
	})

	t.Run("remove -- custom lic ignore not found", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"remove",
			"-l",
			".licignoreNotFound",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "custom ignore file not found: .licignoreNotFound")
	})

	t.Run("remove -- template mit (no error)", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"remove",
			"-t",
			"mit",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.NoError(t, err)
	})

}
