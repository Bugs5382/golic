package internal

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
func (r *ServiceRunner) MustRun() (int, error) {
	log.Info().Msgf("%s command %s started", emoji.Tractor, r.service)
	if err := r.service.Run(); err != nil {
		log.Error().Err(err).Msgf("%s command %s failed", emoji.Bomb, r.service)
		return 1, err
	}
	return r.service.ExitCode(), nil
}
