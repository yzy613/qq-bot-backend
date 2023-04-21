package consts

import "runtime"

const (
	ProjName = "qq-bot-backend"
	Version  = "v1.2.8"
)

var (
	BuildTime   = ""
	CommitHash  = ""
	Description = "Go Version: " + runtime.Version() +
		"\nVersion: " + Version +
		"\nBuild Time: " + BuildTime +
		"\nCommit Hash: " + CommitHash
)
