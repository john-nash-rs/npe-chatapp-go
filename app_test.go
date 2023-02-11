package main

import (
	"fmt"
	"net"
	"runtime"
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
		t.Fatalf(`We have an error. Client Connection is not eastablished`)
	}

	err = client.Send("/topic/whatever", "text/plain", []byte("hello Mr. X! you were awaited."))

	if client == nil || err != nil {
		t.Fatalf(`We have an error. Client Connection is not sending data to a topic`)
	}

	defer client.Disconnect()
}

func TestSendToQueuesAndTopics(t *testing.T) {
	fmt.Println("-------- Started Testing -----")
	ch := make(chan bool, 2)
	fmt.Println("number cpus:", runtime.NumCPU())

	addr := ":59091"

	// channel to communicate that the go routine has started
	started := make(chan bool)

	count := 100
	go runReceiver(ch, count, "/topic/test-1", addr, started)
	<-started
	go runReceiver(ch, count, "/topic/test-1", addr, started)
	<-started
	go runReceiver(ch, count, "/topic/test-2", addr, started)
	<-started
	go runReceiver(ch, count, "/topic/test-2", addr, started)
	<-started
	go runReceiver(ch, count, "/topic/test-1", addr, started)
	<-started
	go runReceiver(ch, count, "/queue/test-1", addr, started)
	<-started
	go runSender(ch, count, "/queue/test-1", addr, started)
	<-started
	go runSender(ch, count, "/queue/test-2", addr, started)
	<-started
	go runReceiver(ch, count, "/queue/test-2", addr, started)
	<-started
	go runSender(ch, count, "/topic/test-1", addr, started)
	<-started
	go runReceiver(ch, count, "/queue/test-3", addr, started)
	<-started
	go runSender(ch, count, "/queue/test-3", addr, started)
	<-started
	go runSender(ch, count, "/queue/test-4", addr, started)
	<-started
	go runSender(ch, count, "/topic/test-2", addr, started)
	<-started
	go runReceiver(ch, count, "/queue/test-4", addr, started)
	<-started

	for i := 0; i < 15; i++ {
		<-ch
	}
}

func runSender(ch chan bool, count int, destination, addr string, started chan bool) {
	conn, _ := net.Dial("tcp", "127.0.0.1"+addr)

	client, _ := stomp.Connect(conn)

	started <- true

	for i := 0; i < count; i++ {
		client.Send(destination, "text/plain",
			[]byte(fmt.Sprintf("%s Hi! Mr X. You were expected. your id is:  %d", destination, i)))
		//println("sent", i)
	}

	ch <- true
}

func runReceiver(ch chan bool, count int, destination, addr string, started chan bool) {
	conn, _ := net.Dial("tcp", "127.0.0.1"+addr)

	client, _ := stomp.Connect(conn)

	sub, _ := client.Subscribe(destination, stomp.AckAuto)

	started <- true

	for i := 0; i < count; i++ {
		msg := <-sub.C
		fmt.Println("----------------   received Message  ------------  ", string(msg.Body))
	}
	ch <- true
}
