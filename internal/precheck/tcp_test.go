package precheck

import (
	"net"
	"testing"
	"time"
)

func TestCheckHost(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	if !CheckHost("127.0.0.1", port, time.Second) {
		t.Fatal("expected host to be reachable")
	}
}

func TestCheckHostInvalid(t *testing.T) {
	if CheckHost("", 0, time.Millisecond) {
		t.Fatal("expected invalid host to fail")
	}
}
