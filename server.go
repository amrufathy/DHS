package main

import (
	"net"
	"fmt"
	"bufio"
)

func main() {
	fmt.Println("Launching server...")

	// create a tcp socket
	listener, _ := net.Listen("tcp", ":5000")
	fmt.Println("Server is listening on port 5000")

	// don't forget to close listener before quitting
	defer listener.Close()

	for {
		// accept a connection
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Connection error %s\n", err)
		}

		handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Print client's remote address
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from", remoteAddr)

	reader := bufio.NewScanner(conn)

	for {
		// get message from client
		recv := reader.Scan()

		if !recv {
			return
		}

		// print message
		fmt.Println("Message Received: ", string(reader.Text()))

		// send reply
		reply := "Hi, Bibo"
		conn.Write([]byte(reply + "\n"))
	}

}
