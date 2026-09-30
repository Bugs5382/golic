package main

/*
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
	"errors"
	"fmt"
	"io"
	"os"

	golog "github.com/Bugs5382/go-log"
	"github.com/Bugs5382/golic/internal"
	"github.com/Bugs5382/golic/internal/commands"
	"github.com/Bugs5382/golic/internal/logging"
	"github.com/enescakir/emoji"
)

// Exit statuses. They are part of the CLI contract: CI jobs run
// `golic inject --dry -x` and treat 1 as "headers missing".
const (
	exitClean   = 0 // nothing to change, or -x not set
	exitChanges = 1 // -x set and files were (or in a dry run would be) modified
	exitError   = 2 // bad flags, bad config, unreadable files and other failures
)

func main() {
	logging.Init(false)
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes golic with args and returns the process exit status.
func run(args []string, stdout, stderr io.Writer) int {
	root := commands.RootCmd()
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	// main reports the error itself, once, with the matching exit status.
	root.SilenceErrors = true

	err := root.Execute()
	if err == nil {
		logging.L().Debug("golic finished clean", golog.F("exit", exitClean))
		return exitClean
	}

	var changes *internal.ChangesError
	if errors.As(err, &changes) {
		logging.L().Debug("golic finished with changes", golog.F("exit", exitChanges), golog.F("files", changes.Count), golog.F("dry", changes.Dry))
		_, _ = fmt.Fprintf(stderr, "%v  %v\n", emoji.Warning, err)
		return exitChanges
	}

	logging.L().Debug("golic failed", golog.F("exit", exitError), golog.F("error", err.Error()))
	_, _ = fmt.Fprintf(stderr, "%v  Error: %v\n", emoji.Bomb, err)
	return exitError
}
