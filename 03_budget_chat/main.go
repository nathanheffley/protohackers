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

		client.Write(fmt.Sprintf("[%s] %s\n", c.Name, message))
	}
}

func (c *Client) Write(message string) {
	fmt.Print("{to ", c.Name, "} ", message)
	c.Conn.Write([]byte(message))
}

func (c *Client) Leave() {
	fmt.Printf("{%s is leaving the room}\n", c.Name)

	if len(clients) < 2 {
		clients = []Client{}
		c.Conn.Close()
		return
	}

	newClients := make([]Client, len(clients)-2)
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

	if len(clients) > 0 {
		clientNames := make([]string, len(clients)-1)
		for _, c := range clients {
			clientNames = append(clientNames, c.Name)
			c.Write("* " + client.Name + " has entered the room\n")
		}
		client.Write(fmt.Sprintf("* The room contains: %s\n", strings.Join(clientNames, ", ")))
	} else {
		client.Write("* The room is empty\n")
	}

	clients = append(clients, client)

	for {
		messageBytes := make([]byte, 1000)
		_, err := conn.Read(messageBytes)
		if err != nil {
			client.Leave()
			break
		}

		message := strings.Trim(string(messageBytes), "\x00")
		message = strings.TrimSpace(message)

		client.SendMessage(string(message))
	}
}
