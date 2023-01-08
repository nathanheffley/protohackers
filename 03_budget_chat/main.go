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
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?"))

	name := make([]byte, 16)
	_, err := conn.Read(name)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	fmt.Println(name)

	conn.Close()
}
