package logging

/*
Apache License 2.0

Copyright 2026 Shane & Contributors

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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildLevels(t *testing.T) {
	t.Run("verbose enables trace output", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("LOG_FORMAT", "")
		var buf bytes.Buffer

		l := build(true, &buf)
		l.Trace().Msg("trace line")
		l.Debug().Msg("debug line")

		assert.Contains(t, buf.String(), "trace line")
		assert.Contains(t, buf.String(), "debug line")
	})

	t.Run("default logs info and above, no trace or debug", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("LOG_FORMAT", "")
		var buf bytes.Buffer

		l := build(false, &buf)
		l.Trace().Msg("trace line")
		l.Debug().Msg("debug line")
		l.Info().Msg("info line")

		assert.NotContains(t, buf.String(), "trace line")
		assert.NotContains(t, buf.String(), "debug line")
		assert.Contains(t, buf.String(), "info line")
	})

	t.Run("LOG_LEVEL wins over verbose", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "warn")
		t.Setenv("LOG_FORMAT", "")
		var buf bytes.Buffer

		l := build(true, &buf)
		l.Trace().Msg("trace line")
		l.Info().Msg("info line")
		l.Warn().Msg("warn line")

		assert.NotContains(t, buf.String(), "trace line")
		assert.NotContains(t, buf.String(), "info line")
		assert.Contains(t, buf.String(), "warn line")
	})

	t.Run("an unknown LOG_LEVEL falls back to the verbose flag", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "loud")
		t.Setenv("LOG_FORMAT", "")
		var buf bytes.Buffer

		l := build(true, &buf)
		l.Trace().Msg("trace line")

		assert.Contains(t, buf.String(), "trace line")
	})
}

func TestBuildFormat(t *testing.T) {
	t.Run("console text by default", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("LOG_FORMAT", "")
		var buf bytes.Buffer

		l := build(false, &buf)
		l.Info().Msg("hello")

		out := buf.String()
		assert.Contains(t, out, "INF")
		assert.Contains(t, out, "hello")
		assert.False(t, strings.HasPrefix(strings.TrimSpace(out), "{"), "default output must not be JSON: %q", out)
	})

	t.Run("JSON when LOG_FORMAT=json", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("LOG_FORMAT", "json")
		var buf bytes.Buffer

		l := build(false, &buf)
		l.Info().Msg("hello")

		var line map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &line))
		assert.Equal(t, "hello", line["message"])
		assert.Equal(t, "info", line["level"])
		assert.Equal(t, "golic", line["service"])
	})
}

func TestInitUnderTestIsSilent(t *testing.T) {
	Init(true)
	assert.False(t, L().Trace().Enabled())
	assert.False(t, L().Error().Enabled())
}
