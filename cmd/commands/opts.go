package commands

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
	Template           string
	ModifiedExitStatus bool
	MasterConfig       string
	Type               LicenseCommandType
}

var InjectOptions Options
var Year = time.Now().Year()
var Company = "MyCompany"
