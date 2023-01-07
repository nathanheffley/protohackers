package main

import (
	"encoding/binary"
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

	ledger := make(map[int]int)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		firstBuf := make([]byte, 4)
		_, err = conn.Read(firstBuf)
		if err != nil {
			panic(err)
		}
		first := int(binary.BigEndian.Uint32(firstBuf))

		secondBuf := make([]byte, 4)
		_, err = conn.Read(secondBuf)
		if err != nil {
			panic(err)
		}
		second := int(binary.BigEndian.Uint32(secondBuf))

		fmt.Println("------")
		fmt.Println(buf[0], first, second)
		fmt.Println("------")

		if buf[0] == 'I' {
			ledger[first] = second
		}

		if buf[0] == 'Q' {
			count := 0
			total := 0
			for timestamp, price := range ledger {
				if timestamp >= first && timestamp <= second {
					count++
					total += price
				}
			}

			if count == 0 {
				binary.BigEndian.PutUint32(output, 0)
			} else {
				binary.BigEndian.PutUint32(output, uint32(total/count))
			}
			break
		}
	}

	conn.Write(output)

	conn.Close()
}
