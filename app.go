package main

import (
	"fmt"
	"net"
	"time"

	"github.com/go-stomp/stomp/server"
)

func main() {
	fmt.Println("Hello, world.")
	addr := "127.0.0.1:59091"
	l, _ := net.Listen("tcp", addr)
	serv := server.Server{
		Addr:          l.Addr().String(),
		Authenticator: nil,
		QueueStorage:  nil,
		HeartBeat:     5 * time.Millisecond,
	}
	serv.Serve(l)
}
