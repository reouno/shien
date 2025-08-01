package version

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func GetVersion() string {
	return Version
}

func GetFullVersion() string {
	if GitCommit == "unknown" {
		return Version
	}
	return Version + "-" + GitCommit[:7]
}