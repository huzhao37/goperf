package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	fPort     = flag.Int("port", 10001, "Port to listen on.")
	fIp       = flag.String("ip", "127.0.0.1", "IP to listen on.")
	keepAlive = flag.Bool("k", true, "is keep alive")
	testType  = flag.Int("test", 0, "test type,ex:0 is concurrency test,1 is IOPS test.")
	buffer    = flag.Int64("b", 500, "send bytes,unit: B.")
)

func receive(conn net.Conn) {
	if !*keepAlive {
		defer conn.Close()
	}
	size := int64(256 * 1024)
	if *testType == 0 {
		size = *buffer
	}
	if *testType == 1 {
		size = 256 * 1024
	}
	buf := make([]byte, size)
	var total uint64
	for {
		n, err := conn.Read(buf)
		total += uint64(n)
		if err != nil {
			fmt.Println("Connection finishes with", total, "bytes:", err)
			return
		}
		if err := binary.Write(conn, binary.BigEndian, total); err != nil {
			panic(err)
		}
	}
}

func accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go receive(conn)
	}
}

func main() {
	flag.Parse()
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.ParseIP(*fIp),
		Port: *fPort,
	})
	if err != nil {
		panic(err)
	}
	go accept(listener)
	//c := make(chan bool)
	//<-c
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
