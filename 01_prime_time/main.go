package main

import (
	"encoding/json"
	"fmt"
	"math"
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
	req := make([]byte, 0)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		if buf[0] == '\n' {
			if err = sendResponse(req, conn); err != nil {
				break
			}
			req = make([]byte, 0)
		} else {
			req = append(req, buf[0])
		}
	}

	conn.Close()
}

func sendResponse(req []byte, conn net.Conn) error {
	if err := ValidateJson(req); err != nil {
		conn.Write([]byte("malformed request"))
		return fmt.Errorf("malformed request")
	}

	var unmarshalled Req
	json.Unmarshal(req, &unmarshalled)

	if checkPrimeNumber(*unmarshalled.Number) {
		conn.Write([]byte("{\"method\":\"isPrime\",\"prime\":true}\n"))
	} else {
		conn.Write([]byte("{\"method\":\"isPrime\",\"prime\":false}\n"))
	}
	return nil
}

type Req struct {
	Method *string
	Number *int
}

func ValidateJson(req []byte) error {
	fmt.Printf("%s\n", req)

	var unmarshalled Req
	if err := json.Unmarshal(req, &unmarshalled); err != nil {
		return fmt.Errorf("invalid JSON")
	}

	if unmarshalled.Method == nil {
		return fmt.Errorf("missing method")
	} else {
		if *unmarshalled.Method != "isPrime" {
			return fmt.Errorf("invalid method")
		}
	}

	if unmarshalled.Number == nil {
		return fmt.Errorf("missing number")
	}

	return nil
}

func checkPrimeNumber(num int) bool {
	if num < 2 {
		return false
	}

	if num == 2 {
		return true
	}

	sq_root := int(math.Sqrt(float64(num)))
	for i := 2; i <= sq_root; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}
