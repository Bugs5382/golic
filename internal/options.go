package internal

/*
Apache License 2.0

Copyright 2006 Shane

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

import "time"

type LicenseCommandType int

const (
	LicenseInject LicenseCommandType = 0
	LicenseRemove LicenseCommandType = 1
)

type Options struct {
	LicIgnore          string
	Copyright          string
	Dry                bool
	ConfigPath         string
	SearchPath         string
	Template           string
	ModifiedExitStatus bool
	MasterConfig       string
	Type               LicenseCommandType
	Verbose            bool
}

var Year = time.Now().Year()

var InjectOptions Options
var RemoveOptions Options
