package internal

import (
	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Run() error
	String() string
	ExitCode() int
}

type ServiceRunner struct {
	service Service
}

// Command Service Runner
func Command(service Service) *ServiceRunner {
	return &ServiceRunner{
		service,
	}
}

// MustRun Run service once and panics if service is broken
func (r *ServiceRunner) MustRun() int {
	log.Info().Msgf("%s command %s started", emoji.Tractor, r.service)
	_ = r.service.Run()
	return r.service.ExitCode()
}
