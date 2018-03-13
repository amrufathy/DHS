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
	"sync"
	"time"
)

var (
	globalNumber = rand.Int63n(100)
	addr         string
	READER       *log.Logger
	WRITER       *log.Logger
	sSequence    int64 = 0
	sNumReaders        = 0
	sNumWriters        = 0
	mutex        sync.Mutex
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
	rLogfile, _ := os.Create("./r_server.log")
	wLogfile, _ := os.Create("./w_server.log")
	READER = log.New(rLogfile, "READ:  ", log.Ltime)
	WRITER = log.New(wLogfile, "WRITE: ", log.Ltime)
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
		recv_text := strings.TrimSpace(scanner.Text())

		if !(len(recv_text) == 0) {
			// print message
			fmt.Println("Message Received:", recv_text)

			if recv_text == "read" {
				sNumReaders++
				serveRead(conn)
			} else if validWriteMessage.Match([]byte(recv_text)) {
				sNumWriters++
				serveWrite(conn, recv_text)
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

	// sleep
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))

	// reply
	err := binary.Write(conn, binary.BigEndian, globalNumber)
	if err != nil {
		fmt.Println(err)
	}

	// log
	READER.Printf("%d\t%d\t%d\t%s\n", sSequence, globalNumber, sNumReaders, conn.RemoteAddr().String())
}

func serveWrite(conn net.Conn, message string) {
	defer decreaseNumWriters()
	defer increaseSequenceNumber()

	// sleep
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))

	// parse message
	newValue, _ := strconv.ParseInt(strings.Fields(message)[1], 10, 64)

	// log
	WRITER.Printf("%d\t%d\t%d\t%s\n", sSequence, globalNumber, newValue, conn.RemoteAddr().String())

	// write newValue
	mutex.Lock() // acquire lock
	globalNumber = newValue
	mutex.Unlock() // release lock
	conn.Write([]byte("valid write\n"))
}

func increaseSequenceNumber() {
	sSequence++
}

func decreaseNumReader() {
	sNumReaders--
}

func decreaseNumWriters() {
	sNumWriters--
}
