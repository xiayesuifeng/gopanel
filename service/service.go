package service

import (
	"context"
	"errors"
	"github.com/coreos/go-systemd/v22/dbus"
	"log"
	"path"
)

type Service struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	ActiveState bool     `json:"activeState"`
	SubStatus   string   `json:"subStatus"`
	TriggeredBy []string `json:"triggeredBy"`
}

type Mode string

const (
	ReplaceMode            = "replace"
	FailMode               = "fail"
	IsolateMode            = "isolate"
	IgnoreDependenciesMode = "ignore-dependencies"
	IgnoreRequirementsMode = "ignore-requirements"
)

var (
	CanceledJobError   = errors.New("systemd service unit job canceled")
	TimeoutJobError    = errors.New("systemd service unit job timeout")
	FailedJobError     = errors.New("systemd service unit job failed")
	DependencyJobError = errors.New("systemd service unit job dependency")
	SkippedJobError    = errors.New("systemd service unit job skipped")
)

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
			triggeredBy := make([]string, 0)
			properties, err := conn.GetUnitPropertiesContext(context, units[0].Name)
			if err != nil {
				log.Println("Failed to get unit properties,unit:", units[0].Name)
			} else {
				if tmp, ok := properties["TriggeredBy"].([]string); ok {
					triggeredBy = tmp
				}
			}

			services = append(services, &Service{
				Name:        units[0].Name,
				Description: units[0].Description,
				Enabled:     file.Type == "enabled",
				ActiveState: units[0].ActiveState == "active",
				SubStatus:   units[0].SubState,
				TriggeredBy: triggeredBy,
			})
		}
	}

	return services, nil
}

func StartService(ctx context.Context, name string, mode Mode) (jobID int, err error) {
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resultChan := make(chan string)
	defer close(resultChan)

	jobID, err = conn.StartUnitContext(ctx, name, string(mode), resultChan)
	if err != nil {
		return
	}

	result := <-resultChan

	switch result {
	case "canceled":
		err = CanceledJobError
	case "timeout":
		err = TimeoutJobError
	case "failed":
		err = FailedJobError
	case "dependency":
		err = DependencyJobError
	case "skipped":
		err = SkippedJobError
	}

	return
}

func StopService(ctx context.Context, name string, mode Mode) (jobID int, err error) {
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resultChan := make(chan string)
	defer close(resultChan)

	jobID, err = conn.StopUnitContext(ctx, name, string(mode), resultChan)

	result := <-resultChan

	switch result {
	case "canceled":
		err = CanceledJobError
	case "timeout":
		err = TimeoutJobError
	case "failed":
		err = FailedJobError
	case "dependency":
		err = DependencyJobError
	case "skipped":
		err = SkippedJobError
	}

	return
}

func RestartService(ctx context.Context, name string, mode Mode) (jobID int, err error) {
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resultChan := make(chan string)
	defer close(resultChan)

	jobID, err = conn.RestartUnitContext(ctx, name, string(mode), resultChan)

	result := <-resultChan

	switch result {
	case "canceled":
		err = CanceledJobError
	case "timeout":
		err = TimeoutJobError
	case "failed":
		err = FailedJobError
	case "dependency":
		err = DependencyJobError
	case "skipped":
		err = SkippedJobError
	}

	return
}
