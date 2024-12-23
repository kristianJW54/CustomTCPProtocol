package main

import (
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	lc := net.ListenConfig{}
	srv := NewServer("test-server", "localhost", "8081", lc)

	go srv.StartServer()
	time.Sleep(1 * time.Second)

	addr := srv.Addr

	if addr == "" {
		t.Error("Failed to start server")
	}

	_, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("Failed to connect to server")
	}
	_, err = net.Dial("tcp", addr)
	if err != nil {
		t.Error("Failed to connect to server")
	}
	_, err = net.Dial("tcp", addr)
	if err != nil {
		t.Error("Failed to connect to server")
	}

	time.Sleep(5 * time.Second)

}
