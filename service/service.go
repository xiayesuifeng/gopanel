package service

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
	"path"
)

type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	ActiveState bool   `json:"activeState"`
	SubStatus   string `json:"subStatus"`
}

func GetServices(context context.Context) ([]*Service, error) {
	conn, err := dbus.NewWithContext(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	unitFiles, err := conn.ListUnitFilesByPatternsContext(context, nil, []string{"*.service"})
	if err != nil {
		return nil, err
	}

	services := make([]*Service, 0)

	for _, file := range unitFiles {
		if file.Type == "static" {
			continue
		}

		units, err := conn.ListUnitsByNamesContext(context, []string{path.Base(file.Path)})
		if err != nil {
			continue
		}

		if len(units) > 0 {
			services = append(services, &Service{
				Name:        units[0].Name,
				Description: units[0].Description,
				Enabled:     file.Type == "enabled",
				ActiveState: units[0].ActiveState == "active",
				SubStatus:   units[0].SubState,
			})
		}
	}

	return services, nil
}
