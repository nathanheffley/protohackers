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

	ledger := make(map[int32]int32)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

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

		secondBuf := make([]byte, 4)
		conn.Read(buf)
		secondBuf[0] = buf[0]
		conn.Read(buf)
		secondBuf[1] = buf[0]
		conn.Read(buf)
		secondBuf[2] = buf[0]
		conn.Read(buf)
		secondBuf[3] = buf[0]

		if messageType == 'I' {
			timestamp := int32(binary.BigEndian.Uint32(firstBuf))
			price := int32(binary.BigEndian.Uint32(secondBuf))
			ledger[timestamp] = price
		}

		if messageType == 'Q' {
			min := int32(binary.BigEndian.Uint32(firstBuf))
			max := int32(binary.BigEndian.Uint32(secondBuf))

			var count, total int32
			count = 0
			total = 0
			for timestamp, price := range ledger {
				if timestamp >= min && timestamp <= max {
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
			conn.Write(output)
			break
		}
	}

	conn.Close()
}
