package firewall

import (
	"errors"
	"gitlab.com/xiayesuifeng/go-firewalld"
	"gitlab.com/xiayesuifeng/gopanel/event"
)

var NotFoundErr = errors.New("not found")

const eventTopic = "firewall"

func GetDefaultZone() (string, error) {
	conn, err := firewalld.New()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.GetDefaultZone()
}

func SetDefaultZone(name string) error {
	conn, err := firewalld.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.SetDefaultZone(name)
	if err != nil {
		return err
	}

	event.Publish(event.Event{
		Topic:   eventTopic,
		Type:    "defaultZoneChanged",
		Payload: name,
	})

	return nil
}

func Reload() error {
	conn, err := firewalld.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Reload()
	if err != nil {
		event.Publish(event.Event{
			Topic:   eventTopic,
			Type:    "reload",
			Payload: nil,
		})
	}

	return nil
}

func Reset() error {
	conn, err := firewalld.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Reset()
	if err != nil {
		event.Publish(event.Event{
			Topic:   eventTopic,
			Type:    "reset",
			Payload: nil,
		})
	}

	return nil
}

func GetICMPTypeNames(permanent bool) ([]string, error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}

	return conn.GetICMPTypeNames()
}

func getConn(permanent bool) (*firewalld.Conn, error) {
	conn, err := firewalld.New()
	if err != nil {
		return nil, err
	}

	if permanent {
		conn = conn.Permanent()
	}

	return conn, nil
}
