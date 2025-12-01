package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"
)

type ChatServer struct {
	mu        sync.Mutex
	clients   []string
	messages  []string
	broadcast chan string
}

func NewChatServer() *ChatServer {
	server := &ChatServer{
		clients:   []string{},
		messages:  []string{},
		broadcast: make(chan string, 10), // channel added
	}
	go server.handleBroadcast()
	return server
}

// Handle sending messages via channel
func (c *ChatServer) handleBroadcast() {
	for msg := range c.broadcast {
		c.mu.Lock()
		c.messages = append(c.messages, msg)
		c.mu.Unlock()
	}
}

func (c *ChatServer) Register(name string, reply *string) error {
	c.mu.Lock()
	c.clients = append(c.clients, name)
	c.mu.Unlock()

	joinMsg := fmt.Sprintf("%s joined the chat", name)
	c.broadcast <- joinMsg
	*reply = "Registered successfully!"
	return nil
}

func (c *ChatServer) SendMessage(msg string, reply *[]string) error {
	c.broadcast <- msg

	c.mu.Lock()
	defer c.mu.Unlock()

	var broadcast []string
	for _, m := range c.messages {
		if !strings.HasPrefix(m, msg[:strings.Index(msg, ":")]) {
			broadcast = append(broadcast, m)
		}
	}
	*reply = broadcast
	return nil
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	listener, _ := net.Listen("tcp", ":1234")
	defer listener.Close()

	fmt.Println("Chat server running on port 1234...")
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}
