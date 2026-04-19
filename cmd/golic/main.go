package main

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
	"fmt"
	"os"

	"github.com/Bugs5382/golic/internal/commands"
	"github.com/Bugs5382/golic/internal/logging"
	"github.com/enescakir/emoji"
)

func main() {

	logging.Init(false)

	if err := commands.RootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v  Error: %v\n", emoji.Bomb, err)
		os.Exit(1)
	}
}
