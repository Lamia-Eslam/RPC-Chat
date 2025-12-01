package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type ChatServer struct {
	mu       sync.Mutex
	messages []string
	users    map[string]bool
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		messages: []string{},
		users:    make(map[string]bool),
	}
}

// Register adds a user and broadcasts join message
func (c *ChatServer) Register(name string, reply *[]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add user if not already registered
	if !c.users[name] {
		c.users[name] = true
		joinMsg := fmt.Sprintf("User %s joined the chat", name)
		c.messages = append(c.messages, joinMsg)
	}

	*reply = c.messages
	return nil
}

// SendMessage appends a new message to the chat
func (c *ChatServer) SendMessage(msg string, reply *[]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = append(c.messages, msg)
	*reply = c.messages
	return nil
}

// GetMessages returns all messages
func (c *ChatServer) GetMessages(_ int, reply *[]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	*reply = c.messages
	return nil
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server running on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
