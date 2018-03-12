package main

import (
	"net"
	"fmt"
	"bufio"
	"github.com/go-ini/ini"
	"strings"
	"strconv"
	"encoding/binary"
	"regexp"
	"math/rand"
	"log"
	"os"
)

var (
	globalNumber = rand.Int63n(100)
	addr         string
	READER       *log.Logger
	WRITER       *log.Logger
	sSequence    int64 = 0
	numReaders   int64 = 0
	numWriters   int64 = 0
)

func init() {
	// This method always executes before main

	// Read configuration file
	cfg, _ := ini.InsensitiveLoad("config.ini")
	info, _ := cfg.GetSection("server")

	host, _ := info.GetKey("host")
	port, _ := info.GetKey("port")
	addr = host.String() + ":" + port.String()

	// Initialize logger
	logfile, _ := os.Create("./server.log")
	READER = log.New(logfile, "READ:  ", log.Ltime)
	WRITER = log.New(logfile, "WRITE: ", log.Ltime)
}

func main() {
	fmt.Println("Launching server...")

	// create a tcp socket
	listener, _ := net.Listen("tcp", addr)
	fmt.Println("Server is listening on", addr)

	// don't forget to close listener before quitting
	defer listener.Close()

	for {
		// accept a connection
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Connection error %s\n", err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Print client's remote address
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from", remoteAddr)
	validWriteMessage := regexp.MustCompile(`write\s[0-9]+`)

	scanner := bufio.NewScanner(conn)

	for {
		// get message from client
		recv := scanner.Scan()

		if !(len(scanner.Text()) == 0) {
			// print message
			fmt.Println("Message Received:", scanner.Text())

			if scanner.Text() == "read" {
				numReaders++
				serveRead(conn)
			} else if validWriteMessage.Match([]byte(scanner.Text())) {
				numWriters++
				serveWrite(conn, scanner.Text())
			} else {
				conn.Write([]byte("Invalid message\n"))
			}
		}

		if !recv {
			break
		}
	}
}

func serveRead(conn net.Conn) {
	defer decreaseNumReader()
	defer increaseSequenceNumber()

	// reply
	err := binary.Write(conn, binary.BigEndian, globalNumber)
	if err != nil {
		fmt.Println(err)
	}

	// log
	READER.Printf("%d\t%d\t%d\t%s\n", sSequence, globalNumber, numReaders, conn.RemoteAddr().String())
}

func serveWrite(conn net.Conn, message string) {
	defer decreaseNumWriters()
	defer increaseSequenceNumber()

	// parse message
	newValue, _ := strconv.ParseInt(strings.Fields(message)[1], 10, 64)

	// log
	WRITER.Printf("%d\t%d\t%d\t%s\n", sSequence, globalNumber, newValue, conn.RemoteAddr().String())

	// write newValue
	globalNumber = newValue
	conn.Write([]byte("valid write\n"))
}

func increaseSequenceNumber() {
	sSequence++
}

func decreaseNumReader() {
	numReaders--
}

func decreaseNumWriters() {
	numWriters--
}
