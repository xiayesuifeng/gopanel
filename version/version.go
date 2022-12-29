package version

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
)

const semverRegexStr = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var (
	// GitCommit git rev-parse HEAD
	GitCommit string

	// Version The main version number that is being run at the moment.
	// Set this variable during `go build` with `-ldflags`:
	//
	//	-ldflags '-X gitlab.com/xiayesuifeng/gopanel/version.Version=x.y.z'
	//
	// for example.
	Version = "1.0.0-dev"

	ErrInvalidSemVer = errors.New("invalid Semantic Versions")
)

type VersionInfo struct {
	Revision   string
	Version    string
	Prerelease string
}

func GetVersion() *VersionInfo {
	reg := regexp.MustCompile(semverRegexStr)

	match := reg.FindStringSubmatch(Version)
	groupNames := reg.SubexpNames()

	if len(match) == 0 {
		log.Panicln(ErrInvalidSemVer)
	}

	result := make(map[string]string)

	for i, name := range groupNames {
		// skip first element
		if i == 0 && name == "" {
			continue
		}

		result[name] = match[i]
	}

	return &VersionInfo{
		Revision:   GitCommit,
		Version:    fmt.Sprintf("%s.%s.%s", result["major"], result["minor"], result["patch"]),
		Prerelease: result["prerelease"],
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
