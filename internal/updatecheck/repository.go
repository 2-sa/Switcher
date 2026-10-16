package updatecheck

import "strings"

const (
	GitHubRepository = "2-sa/Switcher"
	platformsBranch  = "go"
)

// LatestReleaseURL is the browser destination for updates published by the fork.
func LatestReleaseURL() string {
	return "https://github.com/" + GitHubRepository + "/releases/latest"
}

func latestReleaseAPIURL() string {
	return "https://api.github.com/repos/" + GitHubRepository + "/releases/latest"
}

func platformsJSONRawURL(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "Platforms.json"
	}
	return "https://raw.githubusercontent.com/" + GitHubRepository + "/refs/heads/" + platformsBranch + "/" + fileName
}
