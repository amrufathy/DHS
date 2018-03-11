package main

import (
	"net"
	"bufio"
	"os"
	"fmt"
)

func main() {
	fmt.Println("Requesting a connection from server...")

	// send request to server for connection
	conn, _ := net.Dial("tcp", "localhost:5000")

	// create a reader (takes input)
	reader := bufio.NewReader(os.Stdin)

	for {
		// read message from user (input)
		fmt.Print(">: ")
		text, _ := reader.ReadString('\n')

		// write message to connection
		fmt.Fprintf(conn, text+"\n")

		message, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Println("Message from server", message)
	}
}
