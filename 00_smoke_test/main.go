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
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	buf := make([]byte, 1024)

	bufLen, _ := conn.Read(buf)

	fmt.Println(bufLen)
	fmt.Println(buf)

	conn.Write(buf)

	conn.Close()
}
