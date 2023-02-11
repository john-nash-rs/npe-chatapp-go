package main

import (
	"net"
	"testing"
	"time"

	"github.com/go-stomp/stomp"
)

func TestConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:59091")

	if conn == nil || err != nil {
		t.Fatalf(`We have an error. Connection is not eastablished`)
	}
	defer conn.Close()

	client, err := stomp.Connect(conn,
		stomp.ConnOpt.HeartBeat(5*time.Millisecond, 5*time.Millisecond),
	)
	if client == nil || err != nil {
		t.Fatalf(`We have an error. Connection is not eastablished`)
	}
	defer client.Disconnect()
}
