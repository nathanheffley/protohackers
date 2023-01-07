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

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		fmt.Println(buf)
	}

	conn.Write(output)

	conn.Close()
}
