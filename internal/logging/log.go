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
	"time"

	golog "github.com/Bugs5382/go-log"
	"github.com/rs/zerolog"
)

// service is the name go-log stamps on every line.
const service = "golic"

// current holds the logger every package logs through. Init swaps it, and the
// commands run in parallel tests, so it is read and written atomically.
var current atomic.Pointer[zerolog.Logger]

func init() {
	Init(false)
}

// L returns the logger golic logs through.
func L() *zerolog.Logger {
	return current.Load()
}

// Init builds the golic logger on go-log. LOG_LEVEL is read by go-log; when it
// is unset or not a level name, verbose selects trace and the default is info.
// Output goes to stderr so it never mixes with what a command prints on stdout:
// human-readable console text by default, JSON when LOG_FORMAT=json. Under
// `go test` logging is disabled.
func Init(verbose bool) {
	if isTest() {
		l := zerolog.Nop()
		current.Store(&l)
		return
	}
	l := build(verbose, os.Stderr)
	current.Store(&l)
}

// build returns the go-log logger for the given verbosity, writing to out. It
// is split from Init so tests can drive it directly.
func build(verbose bool, out io.Writer) zerolog.Logger {
	l := golog.New(service)
	level := os.Getenv("LOG_LEVEL")

	if _, err := zerolog.ParseLevel(strings.ToLower(level)); level == "" || err != nil {
		if verbose {
			l = l.Level(zerolog.TraceLevel)
		} else {
			l = l.Level(zerolog.InfoLevel)
		}
	}

	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		return l.Output(out)
	}
	return l.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
	})
}

func isTest() bool {
	return strings.HasSuffix(os.Args[0], ".test") ||
		strings.Contains(strings.Join(os.Args, " "), "-test.")
}
