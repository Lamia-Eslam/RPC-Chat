package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var myID int
	client.Call("ChatServer.Register", name, &myID)
	fmt.Printf("Registered successfully! Your ID is %d\n", myID)
	fmt.Println("Start chatting! Type 'exit' to quit.")

	var chatHistory []string
	lastIndex := 0

	// Goroutine to fetch new messages periodically
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			var messages []string
			client.Call("ChatServer.GetMessages", 0, &messages)

			if len(messages) > lastIndex {
				for _, msg := range messages[lastIndex:] {
					// Skip own messages
					if !strings.HasPrefix(msg, fmt.Sprintf("User %d (%s):", myID, name)) {
						fmt.Println("\n" + msg)
					}
				}
				lastIndex = len(messages)
				fmt.Print("> ")
			}
		}
	}()

	// Main loop to send messages
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		fullMsg := fmt.Sprintf("User %d (%s): %s", myID, name, text)
		client.Call("ChatServer.SendMessage", fullMsg, &chatHistory)
		lastIndex = len(chatHistory)
	}
}
