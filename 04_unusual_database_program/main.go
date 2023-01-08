package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.ListenPacket("udp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	db := make(map[string]string)
	db["version"] = "1.0.0"

	for {
		buf := make([]byte, 1024)
		_, addr, err := listener.ReadFrom(buf)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}

		packet := string(buf)

		if strings.Contains(packet, "=") {
			parts := strings.SplitN(packet, "=", 2)
			key := strings.TrimSpace(strings.Trim(parts[0], "\x00"))
			value := strings.TrimSpace(strings.Trim(parts[1], "\x00"))

			if key == "version" {
				fmt.Printf("Attempt to set 'version' to '%s' was ignored\n", value)
				continue
			}

			db[key] = value

			fmt.Printf("SET '%s' => '%s'\n", key, value)

			continue
		}

		key := strings.TrimSpace(strings.Trim(packet, "\x00"))
		value := db[key]
		resp := fmt.Sprintf("%s=%s", key, value)
		listener.WriteTo([]byte(resp), addr)

		fmt.Printf("GET '%s' => '%s'\n", key, value)
	}
}
