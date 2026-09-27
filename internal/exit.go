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

import "fmt"

// ChangesError reports that files were modified, or in a dry run would be,
// while -x (--modified-exit) is set. It is not a failure of golic itself:
// main maps it to exit status 1, and every other error to 2.
type ChangesError struct {
	Type  LicenseCommandType
	Dry   bool
	Count int
}

func (e *ChangesError) Error() string {
	if !e.Dry {
		return fmt.Sprintf("%d file(s) modified", e.Count)
	}
	switch e.Type {
	case LicenseRemove:
		return fmt.Sprintf("%d file(s) have a license header to remove", e.Count)
	case LicenseReplace:
		return fmt.Sprintf("%d file(s) need the license header replaced", e.Count)
	default:
		return fmt.Sprintf("%d file(s) need a license header", e.Count)
	}
}
