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
	MasterConfig       string
	Type               LicenseCommandType
	Verbose            bool
}

var Version = "local" // value can be overridden by ldflags
var Year = time.Now().Year()
var Company = "MyCompany"

var InjectOptions Options
var RemoveOptions Options
