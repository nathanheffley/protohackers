package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	buf := make([]byte, 1)

	ledger := make(map[int32]int32)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(buf)

		messageType := buf[0]

		firstBuf := make([]byte, 4)
		_, err = conn.Read(firstBuf)
		if err != nil {
			panic(err)
		}
		fmt.Println(firstBuf)

		secondBuf := make([]byte, 4)
		_, err = conn.Read(secondBuf)
		if err != nil {
			panic(err)
		}
		fmt.Println(secondBuf)

		if messageType == 'I' {
			// ledger[]
		}
	}

	conn.Write(output)

	conn.Close()
}
