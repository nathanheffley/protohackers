package main

import (
	"encoding/binary"
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
	ledger := make(map[int32]int32)

	output := make([]byte, 4)

	for {
		message, err := io.ReadAll(io.LimitReader(conn, 9))
		if err != nil {
			continue
		}

		fmt.Println(message)

		if len(message) < 1 {
			break
		}

		if message[0] == 'I' {
			timestamp := int32(binary.BigEndian.Uint32(message[1:5]))
			price := int32(binary.BigEndian.Uint32(message[5:]))
			ledger[timestamp] = price
			fmt.Println(timestamp, price)
			continue

		}

		if message[0] == 'Q' {
			min := int32(binary.BigEndian.Uint32(message[1:5]))
			max := int32(binary.BigEndian.Uint32(message[5:]))

			var count, total int64
			count = 0
			total = 0
			for timestamp, price := range ledger {
				if timestamp >= min && timestamp <= max {
					count++
					total += int64(price)
				}
			}

			if count == 0 {
				binary.BigEndian.PutUint32(output, 0)
			} else {
				binary.BigEndian.PutUint32(output, uint32(total/count))
			}
			fmt.Println(output)
			conn.Write(output)
			continue
		}

		fmt.Println("Unexpected Code")
		break
	}

	conn.Close()
}
