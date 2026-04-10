package impl

import (
	"context"

	"github.com/AbsaOSS/golic/internal"
	"github.com/denormal/go-gitignore"
)

type Process struct {
	Ctx  context.Context
	Opts internal.Options

	ignore   gitignore.GitIgnore
	modified bool
}
