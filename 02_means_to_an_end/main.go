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
	// buf := make([]byte, 1)

	// ledger := make(map[int32]int32)

	// output := make([]byte, 4)

	for {
		message, err := io.ReadAll(io.LimitReader(conn, 9))
		if err != nil {
			continue
		}

		fmt.Println(message)

		// if messageType == 'I' {
		// 	timestamp := int32(binary.BigEndian.Uint32(firstBuf))
		// 	price := int32(binary.BigEndian.Uint32(secondBuf))
		// 	ledger[timestamp] = price
		// 	continue
		// }

		// if messageType == 'Q' {
		// 	min := int32(binary.BigEndian.Uint32(firstBuf))
		// 	max := int32(binary.BigEndian.Uint32(secondBuf))

		// 	var count, total int32
		// 	count = 0
		// 	total = 0
		// 	for timestamp, price := range ledger {
		// 		if timestamp >= min && timestamp <= max {
		// 			count++
		// 			total += price
		// 		}
		// 	}

		// 	if count == 0 {
		// 		binary.BigEndian.PutUint32(output, 0)
		// 	} else {
		// 		binary.BigEndian.PutUint32(output, uint32(total/count))
		// 	}
		// 	fmt.Println(output)
		// 	conn.Write(output)
		// 	break
		// }

		fmt.Println("Unexpected Code")
		break
	}

	conn.Close()
}
