package firewall

import (
	"errors"
	"gitlab.com/xiayesuifeng/go-firewalld"
)

var NotFoundErr = errors.New("not found")

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

	return conn.SetDefaultZone(name)
}

func Reload() error {
	conn, err := firewalld.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Reload()
}

func Reset() error {
	conn, err := firewalld.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Reset()
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
