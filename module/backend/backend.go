package backend

import (
	"bytes"
	"io"
	"os/exec"
	"path/filepath"
)

var backendList = map[string]*Backend{}

type Backend struct {
	Status string
	Cmd    *exec.Cmd
	Log    bytes.Buffer
	Notify chan bool

	forceStop bool
}

func StartNewBackend(name, path string, autoReboot bool, arg ...string) {
	b := &Backend{}
	b.Notify = make(chan bool)

	backendList[name] = b

	go func() {
		for !b.forceStop {
			cmd, stdout := b.getCmd(path, arg...)
			b.Cmd = cmd

			b.Log.WriteString("start " + name + "\n")
			go b.readLog(stdout)
			b.SetStatus("Run")
			err := cmd.Run()
			b.Log.WriteString(err.Error() + "\n")
			b.SetStatus("Stop")
			b.Log.WriteString(name + " is dead\n")

			if !autoReboot {
				break
			}
		}
	}()
}

func (b *Backend) getCmd(path string, arg ...string) (*exec.Cmd, io.ReadCloser) {
	cmd := exec.Command(path, arg...)
	cmd.Dir = filepath.Dir(path)

	cmd.Stderr = cmd.Stdout

	stdout, _ := cmd.StdoutPipe()

	return cmd, stdout
}
func (b *Backend) readLog(stdout io.ReadCloser) {
	tmp := make([]byte, 1024)
	for {
		n, err := stdout.Read(tmp)
		b.Log.Write(tmp[:n])
		b.seedNotify()
		if err != nil {
			break
		}
	}
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

func (b *Backend) Stop() error {
	if b.Status == "Run" {
		b.forceStop = true
		return b.Cmd.Process.Kill()
	}

	return nil
}
