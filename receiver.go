package goperf

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
)

var (
	flocalPort = flag.Int("port", 10001, "Port to listen on.")
	flocalIp   = flag.String("ip", "127.0.0.1", "IP to listen on.")
	keepAlive  = flag.Bool("k", false, "is keep alive")
	testType   = flag.Bool("io", false, "test type,either 'concurrency' or 'IOPS'.")
	buffer     = flag.Int64("b", 500, "send bytes,unit: B.")
)

func receive(conn net.Conn) {
	if !*keepAlive {
		defer conn.Close()
	}
	size := int64(256 * 1024)
	if !*testType {
		size = *buffer
	}
	if *testType {
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

func Server() {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.ParseIP(*flocalIp),
		Port: *flocalPort,
	})
	if err != nil {
		panic(err)
	}
	go accept(listener)
}
