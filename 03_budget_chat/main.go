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

func (c *Client) SendMessage(message string) {
	for _, client := range clients {
		if client.Name == c.Name {
			continue
		}

		client.Write(fmt.Sprintf("[%s] %s", c.Name, message))
	}
}

func (c *Client) Write(message string) {
	fmt.Print(message)
	c.Conn.Write([]byte(message))
}

func (c *Client) Leave() {
	fmt.Println("{", c.Name, "is leaving the room}")

	if len(clients) < 2 {
		clients = []Client{}
		c.Conn.Close()
		return
	}

	newClients := make([]Client, len(clients)-1)
	for _, client := range clients {
		if client.Name != c.Name {
			newClients = append(newClients, client)
			client.Write(fmt.Sprintf("* %s has left the room\n", c.Name))
		}
	}
	clients = newClients
	c.Conn.Close()
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

	nameBytes := make([]byte, 16)
	_, err := conn.Read(nameBytes)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	name := strings.Trim(string(nameBytes), "\x00")
	name = strings.TrimSpace(name)

	if len(name) < 1 {
		conn.Close()
		return
	}

	client := Client{
		Name: string(name),
		Conn: conn,
	}

	clientNames := make([]string, len(clients))
	for _, c := range clients {
		clientNames = append(clientNames, c.Name)
		c.Write("* " + client.Name + " has entered the room\n")
	}
	client.Write(fmt.Sprintf("* The room contains: %s\n", strings.Join(clientNames, ", ")))

	clients = append(clients, client)

	for {
		message := make([]byte, 1000)
		_, err := conn.Read(message)
		if err != nil {
			client.Leave()
			break
		}

		client.SendMessage(string(message))
	}

	client.Leave()
}
