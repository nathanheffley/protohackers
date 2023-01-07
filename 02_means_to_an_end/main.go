package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
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
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 9)

	ledger := make(map[int]int)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		fmt.Println(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(buf[1:5])
		first := int(binary.BigEndian.Uint32(buf[1:5]))
		second := int(binary.BigEndian.Uint32(buf[5:9]))

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
			fmt.Println(output)
			break
		}
	}

	conn.Write(output)

	conn.Close()
}
