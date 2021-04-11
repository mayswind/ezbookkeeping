package version

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var (
	// Version holds the version of this execution program
	Version string

	// CommitHash holds the git commit hash of this execution program's source code
	CommitHash string

	// BuildUnixTime holds the time when starting building this execution program
	BuildUnixTime string
)

// GetFullVersion returns the full version
func GetFullVersion() string {
	fullVersion := "Local Build"

	if Version != "" {
		fullVersion = Version
	}

	additionalInfos := make([]string, 0, 2)

	if CommitHash != "" {
		additionalInfos = append(additionalInfos, "commit "+CommitHash)
	}

	if BuildUnixTime != "" {
		unixTime, err := utils.StringToInt64(BuildUnixTime)

		if unixTime > 0 && err == nil {
			additionalInfos = append(additionalInfos, "build time "+utils.FormatUnixTimeToLongDateTimeInServerTimezone(unixTime))
		}
	}

	if len(additionalInfos) > 0 {
		fullVersion = fmt.Sprintf("%s (%s)", fullVersion, strings.Join(additionalInfos, ", "))
	}

	return fullVersion
}
