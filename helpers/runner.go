package helpers

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
	_ = r.service.Run()
	return r.service.ExitCode()
}
