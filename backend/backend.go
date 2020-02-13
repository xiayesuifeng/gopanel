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
	Notify chan bool
}

func StartNewBackend(name, path string, arg ...string) {
	cmd := exec.Command(path, arg...)
	cmd.Dir = filepath.Dir(path)
	b := &Backend{Cmd: cmd}
	b.Notify = make(chan bool)

	cmd.Stderr = cmd.Stdout

	stdout, _ := cmd.StdoutPipe()

	backendList[name] = b

	go func() {
		b.SetStatus("Run")
		cmd.Run()
		b.SetStatus("Stop")
	}()

	go func() {
		tmp := make([]byte, 1024)
		for {
			if b.Status != "Run" {
				continue
			}
			n, _ := stdout.Read(tmp)
			b.Log.Write(tmp[:n])
			b.seedNotify()
		}
	}()
}

func (b *Backend) seedNotify() {
	select {
	case b.Notify <- true:
	default:
	}
}

func GetBackend(name string) *Backend {
	return backendList[name]
}

func (b *Backend) SetStatus(status string) {
	b.Status = status
	b.seedNotify()
}
