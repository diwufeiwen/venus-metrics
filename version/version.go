package version

var CurrentCommit string

// BuildVersion is the local build version, set by build system
const BuildVersion = "v0.1.0"

var UserVersion = BuildVersion + CurrentCommit
