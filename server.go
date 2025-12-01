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
	users    map[int]string // ID -> Name
	nextID   int
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		messages: []string{},
		users:    make(map[int]string),
		nextID:   1,
	}
}

// Register a client and broadcast join message
func (c *ChatServer) Register(name string, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID
	c.nextID++
	c.users[id] = name

	joinMsg := fmt.Sprintf("User %d (%s) joined", id, name)
	c.messages = append(c.messages, joinMsg)

	*reply = id
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
