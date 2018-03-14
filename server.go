package main

import (
	"net"
	"fmt"
	"bufio"
	"github.com/go-ini/ini"
	"strings"
	"strconv"
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
	rLogfile, _ := os.Create("./logs/server/r.log")
	wLogfile, _ := os.Create("./logs/server/w.log")
	READER = log.New(rLogfile, "READ:  ", log.Ltime)
	WRITER = log.New(wLogfile, "WRITE: ", log.Ltime)
}

func main() {
	fmt.Println("Launching server...")

	// create a tcp socket
	addr, _ := net.ResolveTCPAddr("tcp", addr)
	listener, _ := net.ListenTCP("tcp", addr)
	fmt.Println("Server is listening on", addr)

	// don't forget to close listener before quitting
	defer listener.Close()

	for {
		// accept a connection
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Printf("Connection error %s\n", err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Print client's remote address
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from", remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		// get message from client
		//var recv_text string
		//if line, err := reader.ReadString('\n'); err != nil {
		//	println("Error: ", err.Error())
		//	break
		//} else {
		//	line = strings.TrimSpace(line)
		//	recv_text = line
		//}
		recv := scanner.Scan()
		recv_text := strings.TrimSpace(scanner.Text())

		if !(len(recv_text) == 0) {
			// print message
			fmt.Println("Message Received:", recv_text)

			if recv_text == "read" {
				sNumReaders++
				serveRead(conn)
			} else if recv_text == "write" {
				sNumWriters++
				serveWrite(conn, *scanner)
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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// reply
	// send data
	writer.WriteString(strconv.FormatInt(globalNumber, 10) + "\n")
	writer.Flush()

	// send sequence number
	writer.WriteString(strconv.FormatInt(sSequence, 10) + "\n")
	writer.Flush()

	var rIdx int64
	if line, err := reader.ReadString('\n'); err != nil {
		println("Error: ", err.Error())
	} else {
		val, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		rIdx = val
	}

	// log
	READER.Printf("%d\t%d\t%d\t%d\t%s\n", sSequence, globalNumber, sNumReaders, rIdx, conn.RemoteAddr().String())
}

func serveWrite(conn net.Conn, scanner bufio.Scanner) {
	defer decreaseNumWriters()
	defer increaseSequenceNumber()

	//scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	var newValue int64

	// read new value
	if recv := scanner.Scan(); !recv {
		println("Error:", scanner.Err().Error())
	} else {
		println(scanner.Text())
		val, _ := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
		newValue = val
	}

	// sleep
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	// send sequence number
	writer.WriteString(strconv.FormatInt(sSequence, 10) + "\n")
	writer.Flush()

	// log
	WRITER.Printf("%d\t%d\t%d\t%s\n", sSequence, globalNumber, newValue, conn.RemoteAddr().String())

	// write newValue
	mutex.Lock() // acquire lock
	globalNumber = newValue
	mutex.Unlock() // release lock
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

func s_read(reader bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			fmt.Println(line)
		} else {
			break
		}
	}
}
