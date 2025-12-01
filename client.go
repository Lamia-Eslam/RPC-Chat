package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

func main() {
	client, _ := rpc.Dial("tcp", "localhost:1234")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var reply string
	client.Call("ChatServer.Register", name, &reply)
	fmt.Println(reply)

	fmt.Println("\nStart chatting! Type 'exit' to quit.")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		fullMsg := fmt.Sprintf("%s: %s", name, text)
		var history []string
		client.Call("ChatServer.SendMessage", fullMsg, &history)

		fmt.Println("\n--- Chat Broadcast ---")
		for _, msg := range history {
			fmt.Println(msg)
		}
		fmt.Println("----------------------")
	}
}
