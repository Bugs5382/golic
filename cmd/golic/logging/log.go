package logging

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
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init sets up logging configuration
func Init(verbose bool) {
	// If we are in a test, stop everything immediately.
	if isTest() {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return
	}

	setLogFormat()

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr != "" {
		if parsedLevel, err := zerolog.ParseLevel(strings.ToLower(levelStr)); err == nil {
			zerolog.SetGlobalLevel(parsedLevel)
			return
		}
	}

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func isTest() bool {
	return strings.HasSuffix(os.Args[0], ".test") ||
		strings.Contains(strings.Join(os.Args, " "), "-test.")
}

func setLogFormat() {
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))

	if format == "json" {
		return
	}

	// Otherwise, default to ConsoleWriter (Text)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
}
