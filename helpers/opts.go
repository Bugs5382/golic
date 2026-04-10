package helpers

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
	Type               LicenseCommandType
	Verbose            bool
}

var InjectOptions Options
var RemoveOptions Options
var Year = time.Now().Year()
var Company = "MyCompany"
