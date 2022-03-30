package caddymodule

import (
	"os/exec"
	"strings"
)

type ModuleList struct {
	Standard    map[string][]string
	NonStandard map[string][]string
	Unknown     map[string][]string
}

func GetModuleList() (*ModuleList, error) {
	path, err := exec.LookPath("caddy")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, "list-modules", "-packages")

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	list := &ModuleList{}

	tmp := make(map[string][]string)
	for _, line := range strings.Split(string(outBytes), "\n") {
		if line == "" {
			continue
		}

		if strings.Contains(line, "modules:") {
			if strings.Contains(line, "Standard") {
				list.Standard = tmp
			} else if strings.Contains(line, "Non-standard") {
				list.NonStandard = tmp
			} else if strings.Contains(line, "Unknown") {
				list.Unknown = tmp
			}

			tmp = make(map[string][]string)
		} else if item := strings.Split(line, " "); len(item) == 2 {
			pkgName := item[1]
			tmp[pkgName] = append(tmp[pkgName], item[0])
		}
	}

	return list, nil
}
