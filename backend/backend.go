package backend

import (
	"bytes"
	"os/exec"
	"path/filepath"
)

var backendList = map[string]*Backend{}

type Backend struct {
	Status string
	Cmd    *exec.Cmd
	Log    bytes.Buffer
}

func StartNewBackend(name, path string, arg ...string) {
	cmd := exec.Command(path, arg...)
	cmd.Dir = filepath.Dir(path)
	b := &Backend{Cmd: cmd}

	cmd.Stdout = &b.Log
	cmd.Stderr = &b.Log

	backendList[name] = b

	go func() {
		b.Status = "Run"
		cmd.Run()
		b.Status = "Stop"
	}()
}

func GetBackend(name string) *Backend {
	return backendList[name]
}
