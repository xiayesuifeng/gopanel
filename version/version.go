package version

import (
	"bytes"
	"fmt"
)

var (
	// GitCommit git rev-parse HEAD
	GitCommit string

	// Version The main version number that is being run at the moment.
	// Set this variable during `go build` with `-ldflags`:
	//
	//	-ldflags '-X gitlab.com/xiayesuifeng/gopanel/version.Version=x.y.z'
	//
	// for example.
	Version string

	// Prerelease A pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	// Set this variable during `go build` with `-ldflags`:
	//
	//	-ldflags '-X gitlab.com/xiayesuifeng/gopanel/version.Prerelease=dev'
	//
	// for example.
	Prerelease = "dev"
)

type VersionInfo struct {
	Revision   string
	Version    string
	Prerelease string
}

func GetVersion() *VersionInfo {
	ver := Version
	rel := Prerelease

	return &VersionInfo{
		Revision:   GitCommit,
		Version:    ver,
		Prerelease: rel,
	}
}

func (v *VersionInfo) VersionNumber() string {
	version := v.Version

	if v.Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, v.Prerelease)
	}

	return version
}

func (v *VersionInfo) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "Gopanel v%s", v.Version)
	if v.Prerelease != "" {
		fmt.Fprintf(&versionString, "-%s", v.Prerelease)
	}

	if rev && v.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", v.Revision)
	}

	return versionString.String()
}
