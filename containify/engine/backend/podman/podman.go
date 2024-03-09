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

	conn, err := bindings.NewConnection(context.Background(), p.uri)
	if err != nil {
		return err
	}

	p.serviceVersion = bindings.ServiceVersion(conn)

	return nil
}

func (p *Podman) Container() engine.Container {
	//TODO implement me
	panic("implement me")
}

func (p *Podman) Image() engine.Image {
	//TODO implement me
	panic("implement me")
}

func init() {
	engine.Register("podman", &Podman{})
}

func (p *Podman) getConn(ctx context.Context) (context.Context, error) {
	return bindings.NewConnection(ctx, p.uri)
}
