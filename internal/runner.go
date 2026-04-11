package internal

import (
	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
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
	log.Infof("%s command %s started", emoji.Tractor, r.service)
	_ = r.service.Run()
	return r.service.ExitCode()
}
