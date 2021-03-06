package app

import (
	"runtime"
	"strings"
)

var (
	Name           = "watermint toolbox"
	Version        = "`dev`"
	Copyright      = "© 2016-2020 Takayuki Okazaki"
	Hash           = ""
	Branch         = ""
	Zap            = ""
	BuilderKey     = ""
	BuildTimestamp = ""
	DefaultWebPort = 7800
)

func UserAgent() string {
	return strings.ReplaceAll(Name, " ", "-") + "/" + Version
}

func ReleaseStage() string {
	switch Branch {
	case "current":
		return StageBeta
	case "master":
		return StageRelease
	default:
		return StageDev
	}
}

func IsProduction() bool {
	return Hash != ""
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
