package firewall

func GetServiceNames(permanent bool) ([]string, error) {
	conn, err := getConn(permanent)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.GetServiceNames()
}
