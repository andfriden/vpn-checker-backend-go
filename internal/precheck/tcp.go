package precheck

import (
	"net"
	"strconv"
	"time"
)

func CheckHost(
	host string,
	port int,
	timeout time.Duration,
) bool {

	if host == "" || port <= 0 {
		return false
	}

	address := net.JoinHostPort(
		host,
		strconv.Itoa(port),
	)

	conn, err := net.DialTimeout(
		"tcp",
		address,
		timeout,
	)

	if err != nil {
		return false
	}

	_ = conn.Close()

	return true
}
