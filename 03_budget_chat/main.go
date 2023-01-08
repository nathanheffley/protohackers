package main

import (
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Name string
	Conn net.Conn
}

func (c *Client) Write(message string) {
	c.Conn.Write([]byte(message + "\n"))
}

func (c *Client) Leave() {
	c.Conn.Close()
	newClients := make([]Client, len(clients)-1)
	for _, client := range clients {
		if client.Name != c.Name {
			newClients = append(newClients, client)
			client.Write(fmt.Sprintf("* %s has left the room", c.Name))
		}
	}
	clients = newClients
}

var clients = []Client{}

func main() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	name := make([]byte, 16)
	_, err := conn.Read(name)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	client := Client{
		Name: string(name),
		Conn: conn,
	}

	clientNames := make([]string, len(clients))
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
		client.Write("* " + client.Name + " has entered the room")
	}
	roomContainsMessage := fmt.Sprintf("* The room contains: %s\n", strings.Join(clientNames, ", "))
	client.Write(roomContainsMessage)
	fmt.Println(roomContainsMessage)

	clients = append(clients, client)

	for {
		message := make([]byte, 1000)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, c := range clients {
			if c.Name == client.Name {
				continue
			}

			c.Write(fmt.Sprintf("[%s] %s", client.Name, string(message)))
		}
	}

	client.Leave()
}
