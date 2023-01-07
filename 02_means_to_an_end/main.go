package main

import (
	"bytes"
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
	buf := make([]byte, 9)

	ledger := make(map[int]int)

	output := make([]byte, 4)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		fmt.Printf("char: %b\n", buf[0])

		reader := bytes.NewBuffer(buf[1:5])
		firstData := make([]byte, 4)
		err = binary.Read(reader, binary.BigEndian, &firstData)
		if err != nil {
			panic(err)
		}
		first := int(binary.BigEndian.Uint32(firstData))
		fmt.Printf("first: %b\n", firstData)

		reader = bytes.NewBuffer(buf[5:9])
		secondData := make([]byte, 4)
		err = binary.Read(reader, binary.BigEndian, &secondData)
		if err != nil {
			panic(err)
		}
		second := int(binary.BigEndian.Uint32(secondData))
		fmt.Printf("second: %b\n", secondData)

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
