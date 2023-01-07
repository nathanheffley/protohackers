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

	// ledger := make(map[int32]int32)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("char", buf)

		messageType := buf[0]

		firstBuf := make([]byte, 4)
		conn.Read(buf)
		firstBuf[0] = buf[0]
		conn.Read(buf)
		firstBuf[1] = buf[0]
		conn.Read(buf)
		firstBuf[2] = buf[0]
		conn.Read(buf)
		firstBuf[3] = buf[0]
		fmt.Println(firstBuf)

		secondBuf := make([]byte, 4)
		conn.Read(buf)
		secondBuf[0] = buf[0]
		conn.Read(buf)
		secondBuf[1] = buf[0]
		conn.Read(buf)
		secondBuf[2] = buf[0]
		conn.Read(buf)
		secondBuf[3] = buf[0]
		fmt.Println(secondBuf)

		if messageType == 'I' {
			// ledger[]
		}
	}

	conn.Write(output)

	conn.Close()
}
