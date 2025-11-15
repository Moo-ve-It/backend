package vcs

import (
	"fmt"
	"runtime/debug"
)

func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range buildInfo.Settings {
			switch s.Key {
			case "vcs.time":
				time = s.Value
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	// If the code is modified, the version will have a "-dirty" suffix
	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	// Otherwise we return the time and commit hash as a version #
	return fmt.Sprintf("%s-%s", time, revision)
}
