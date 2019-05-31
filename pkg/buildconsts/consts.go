package buildconsts

import "fmt"

var (
	// BuildDate is the date and time at which the platform was built.
	BuildDate string

	// GitCommit is the SHA of the commit from which the platform was built.
	GitCommit string

	// GitBranch is the branch name from which the platform was built.
	GitBranch string

	// GitState indicates whether uncommitted changes were present at build
	// time.
	GitState string

	// GitSummary is the output of `git describe --tags --dirty --always` at
	// build time.
	GitSummary string

	// ProjectName is the name of the complete system, used for config
	// directories etc
	ProjectName = "nomad-toast"

	// Version of the build
	Version string
)

// GetVersion is responsible for populating and returning the binary version information.
func GetVersion() string {
	var gitCommit string
	if len(GitCommit) >= 8 {
		gitCommit = GitCommit[:7]
	}

	return fmt.Sprintf("%s %s\n Date: %s\n Commit: %s\n Branch: %s\n State: %s\n Summary:%s",
		ProjectName, Version, BuildDate, gitCommit, GitBranch, GitState, GitSummary)
}
