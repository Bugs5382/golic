package internal

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
