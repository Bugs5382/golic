package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/AbsaOSS/golic/cmd/commands"
	"github.com/AbsaOSS/golic/cmd/logging"
	"github.com/enescakir/emoji"
)

//go:embed .golic.yaml
var golicConfig string

func main() {

	logging.Init()

	if err := commands.RootCmd(golicConfig).Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v  Error: %v\n", emoji.Bomb, err)
		os.Exit(1)
	}
}
