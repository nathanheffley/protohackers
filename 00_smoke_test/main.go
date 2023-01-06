package main

import "net"

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	buf := make([]byte, 1024)

	conn.Read(buf)

	conn.Write(buf)

	conn.Close()
}
