package main

import (
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

	for {
		bufLen, err := conn.Read(buf)
		if err != nil {
			if bufLen > 0 {
				conn.Write(buf[0:bufLen])
			}
			break
		}

		if bufLen > 0 {
			conn.Write(buf[0:bufLen])
		}
	}

	conn.Close()
}
