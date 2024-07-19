package firewall

import "gitlab.com/xiayesuifeng/go-firewalld"

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
