package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const maxMessageLength = 1000

type client struct {
	name   string
	conn   net.Conn
	reader *bufio.Reader
}

func (c *client) read() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println(c.name, "disconnected")
			c.conn.Close()
			return
		}

		if len(line) > maxMessageLength {
			fmt.Println(c.name, "sent a message that exceeded the maximum length")
			continue
		}

		if line[0] == '/' {
			parts := strings.Split(line, " ")
			if parts[0] == "/setname" && len(parts) > 1 {
				c.name = parts[1]
				fmt.Println(c.name, "has joined the chat room")
				sendToAllClients(fmt.Sprintf("* %s has entered the room", c.name))
				continue
			}
		}

		sendToAllClients(fmt.Sprintf("[%s] %s", c.name, line))
	}
}

func (c *client) write() {
	for {
		message, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println(c.name, "disconnected")
			c.conn.Close()
			return
		}

		if len(message) > maxMessageLength {
			fmt.Println(c.name, "sent a message that exceeded the maximum length")
			continue
		}

		sendToAllClients(fmt.Sprintf("[%s] %s", c.name, message))
	}
}

func sendToAllClients(message string) {
	for _, c := range clients {
		c.conn.Write([]byte(message))
	}
}

var clients []*client

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Listening for connections...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		reader := bufio.NewReader(conn)
		client := &client{
			conn:   conn,
			reader: reader,
			name:   "anonymous",
		}
		clients = append(clients, client)
		go client.read()
	}
}
