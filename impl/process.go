package impl

import (
	"context"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/denormal/go-gitignore"
)

type Process struct {
	Ctx  context.Context
	Opts helpers.Options

	ignore   gitignore.GitIgnore
	modified bool
}
