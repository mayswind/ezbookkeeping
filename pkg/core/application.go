package core

import "fmt"

// ApplicationName represents the application name
const ApplicationName = "ezBookkeeping"

// Version, CommitHash and BuildTime are set at build
var (
	Version    string
	CommitHash string
	BuildTime  string
)

func GetOutgoingUserAgent() string {
	if Version == "" {
		return ApplicationName
	}

	return fmt.Sprintf("%s/%s", ApplicationName, Version)
}
