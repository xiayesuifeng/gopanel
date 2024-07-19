package firewall

import "gitlab.com/xiayesuifeng/go-firewalld"

func GetDefaultZone() (string, error) {
	conn, err := firewalld.New()
	if err != nil {
		return "", err
	}

	return conn.GetDefaultZone()
}
