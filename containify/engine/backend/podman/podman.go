package podman

import (
	"context"
	"encoding/json"
	"github.com/blang/semver/v4"
	"github.com/containers/podman/v5/pkg/bindings"
	"gitlab.com/xiayesuifeng/gopanel/containify/engine"
)

type Podman struct {
	uri string

	serviceVersion *semver.Version

	container *container
	image     *image
}

type Setting struct {
	// Endpoint podman api socket uri
	Endpoint string `json:"endpoint"`
}

func (p *Podman) New(setting []byte) error {
	data := &Setting{}

	if err := json.Unmarshal(setting, data); err != nil {
		return err
	}

	p.uri = data.Endpoint

	//p.uri = "unix:///run/user/1000/podman/podman.sock"
	p.uri = "unix:///var/run/podman/podman.sock"
	conn, err := bindings.NewConnection(context.Background(), p.uri)
	if err != nil {
		return err
	}

	p.serviceVersion = bindings.ServiceVersion(conn)

	p.container = &container{podman: p}
	p.image = &image{podman: p}

	return nil
}

func (p *Podman) Container() engine.Container {
	return p.container
}

func (p *Podman) Image() engine.Image {
	return p.image
}

func init() {
	engine.Register("podman", &Podman{})
}

func (p *Podman) getConn(ctx context.Context) (context.Context, error) {
	return bindings.NewConnection(ctx, p.uri)
}
