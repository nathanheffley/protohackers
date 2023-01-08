package main

import (
	"fmt"
	"io"
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
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?"))

	name, err := io.ReadAll(io.LimitReader(conn, 9))
	if err != nil {
		conn.Close()
		return
	}

	fmt.Println(name)

	conn.Close()
}
