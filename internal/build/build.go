package build

var (
	version = "vSNAPSHOT"
	time    = "unknown"
	branch  = "n/a"
	commit  = "n/a"
)

func Version() string {
	return version
}

func Time() string {
	return time
}

func Branch() string {
	return branch
}

func Commit() string {
	return commit
}
