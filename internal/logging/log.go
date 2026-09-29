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
	"io"
	"os"
	"strings"
	"sync/atomic"

	golog "github.com/Bugs5382/go-log"
)

// service is the name go-log stamps on every line.
const service = "golic"

var (
	// current holds the logger every package logs through. Init swaps it, and
	// the commands run in parallel tests, so it is read and written atomically.
	current atomic.Pointer[golog.TraceLogger]
	// debug records whether the current logger writes debug lines, since the
	// neutral go-log API does not expose its level.
	debug atomic.Bool
)

func init() {
	Init(false)
}

// L returns the logger golic logs through.
func L() golog.TraceLogger {
	return *current.Load()
}

// DebugEnabled reports whether the current logger writes debug lines, so a
// caller can skip building output that would only be thrown away.
func DebugEnabled() bool {
	return debug.Load()
}

// Init builds the golic logger on go-log. A valid LOG_LEVEL wins; when it is
// unset or not a level name, verbose selects trace and the default is info.
// Output goes to stderr so it never mixes with what a command prints on stdout:
// human-readable console text by default, JSON when LOG_FORMAT=json. Under
// `go test` logging is disabled.
func Init(verbose bool) {
	l := golog.Nop()
	on := false
	if !isTest() {
		l = build(verbose, os.Stderr)
		on = debugOn(verbose)
	}
	current.Store(&l)
	debug.Store(on)
}

// build returns the go-log logger for the given verbosity, writing to out. It
// is split from Init so tests can drive it directly.
func build(verbose bool, out io.Writer) golog.TraceLogger {
	return golog.NewLoggerWithOptions(service,
		golog.WithOutput(out),
		golog.WithDefaultFormat(golog.FormatConsole),
		golog.WithDefaultLevel(defaultLevel(verbose)),
	)
}

func defaultLevel(verbose bool) golog.Level {
	if verbose {
		return golog.LevelTrace
	}
	return golog.LevelInfo
}

// debugOn mirrors how go-log picks the level: LOG_LEVEL when it names a
// level, the verbose default otherwise.
func debugOn(verbose bool) bool {
	level := defaultLevel(verbose)
	switch env := golog.Level(strings.ToLower(os.Getenv("LOG_LEVEL"))); env {
	case golog.LevelTrace, golog.LevelDebug, golog.LevelInfo, golog.LevelWarn,
		golog.LevelError, golog.LevelFatal, golog.LevelPanic, golog.LevelDisabled:
		level = env
	}
	return level == golog.LevelTrace || level == golog.LevelDebug
}

func isTest() bool {
	return strings.HasSuffix(os.Args[0], ".test") ||
		strings.Contains(strings.Join(os.Args, " "), "-test.")
}
