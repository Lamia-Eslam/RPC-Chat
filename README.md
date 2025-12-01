# Go RPC Chat App
This is a simple real-time chat program in Go. Multiple clients can connect to a server, send messages, and see messages from other users. Your own messages are not shown back to you.

How to Run ?

# Start the Server
Open a terminal in the project folder and run:
go run server.go

# Start the Clients
Open one or more terminals for the clients and run:
go run client.go
-----------------------------------------------------
# Key Features:
- Multiple clients support – several clients can connect and chat at the same time.
- Real-time broadcasting – messages from one client appear to all other clients immediately.
- Join notifications – all clients are notified when a new user joins.
- No self-echo – a client does not see its own messages.
- Concurrency handling – server uses mutex and clients use goroutines to safely manage shared resources.
