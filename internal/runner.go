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
	"github.com/Bugs5382/golic/internal/logging"
	"github.com/enescakir/emoji"
)

type Service interface {
	Run() error
	String() string
	// Changes reports files that were or would be modified when the caller
	// asked for that to count as a non-clean result (-x).
	Changes() error
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

// MustRun runs the service once. It returns the run error, or the
// *ChangesError from Changes when the run succeeded but files need a change.
func (r *ServiceRunner) MustRun() error {
	logging.L().Info().Msgf("%s command %s started", emoji.Tractor, r.service)
	if err := r.service.Run(); err != nil {
		logging.L().Error().Err(err).Msgf("%s command %s failed", emoji.Bomb, r.service)
		return err
	}
	if err := r.service.Changes(); err != nil {
		logging.L().Debug().Err(err).Msgf("command %s finished with changes", r.service)
		return err
	}
	logging.L().Debug().Msgf("command %s finished clean", r.service)
	return nil
}
