package consts

import (
	"github.com/gogf/gf/v2"
	"runtime"
)

const (
	ProjName = "qq-bot-backend"
	Version  = "1.8.3"
)

var (
	GitTag      = ""
	GitCommit   = ""
	BuildTime   = ""
	Description = "Version: " + Version +
		"\nGo Version: " + runtime.Version() +
		"\nGoFrame Version: " + gf.VERSION +
		"\nGit Tag: " + GitTag +
		"\nGit Commit: " + GitCommit +
		"\nBuild Time: " + BuildTime
)
