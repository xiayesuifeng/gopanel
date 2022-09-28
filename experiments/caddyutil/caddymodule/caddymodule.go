package caddymodule

import (
	"os/exec"
	"strings"
)

type ModuleList struct {
	Standard    map[string][]string `json:"standard"`
	NonStandard map[string][]string `json:"nonStandard"`
	Unknown     map[string][]string `json:"unknown"`
}

func GetModuleList() (*ModuleList, error) {
	path, err := exec.LookPath("caddy")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, "list-modules", "--packages")

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

func (m *ModuleList) HasNonStandardModule(module string) bool {
	for _, list := range m.NonStandard {
		for _, name := range list {
			if name == module {
				return true
			}
		}
	}

	return false
}

func (m *ModuleList) HasPackage(pkg string) bool {
	_, ok := m.Standard[pkg]
	if ok {
		return true
	}

	_, ok = m.NonStandard[pkg]
	if ok {
		return true
	}

	_, ok = m.Unknown[pkg]
	if ok {
		return true
	}

	return false
}
