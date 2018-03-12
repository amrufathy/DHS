package main

import (
	"net"
	"bufio"
	"os"
	"fmt"
	"github.com/go-ini/ini"
	"encoding/binary"
	"strings"
	"strconv"
)

var (
	sAddr       string
	cNumReaders int
	cNumWriters int
	READ_MSG    = "read"
	WRITE_MSG   = "write %d"
)

func init() {
	// Read configuration file
	cfg, _ := ini.InsensitiveLoad("config.ini")
	sInfo, _ := cfg.GetSection("server")

	sHost, _ := sInfo.GetKey("host")
	sPort, _ := sInfo.GetKey("port")
	sAddr = sHost.String() + ":" + sPort.String()

	cInfo, _ := cfg.GetSection("client")
	cReaders, _ := cInfo.GetKey("num_readers")
	cWriters, _ := cInfo.GetKey("num_writers")

	cReadersInt, _ := strconv.Atoi(cReaders.String())
	cWritersInt, _ := strconv.Atoi(cWriters.String())

	cNumReaders = cReadersInt
	cNumWriters = cWritersInt
}

func main() {
	fmt.Println("Requesting a connection from server...")

	// send request to server for connection
	conn, _ := net.Dial("tcp", sAddr)

	// create a reader (takes input)
	reader := bufio.NewReader(os.Stdin)

	for {
		// read message from user (input)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		// write message to connection
		fmt.Fprintf(conn, text)

		if strings.Contains(text, "read") {
			var number int64
			err := binary.Read(conn, binary.BigEndian, &number)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Message from server", number)
		} else if strings.Contains(text, "write") {
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println("Message from server", message)
		}
	}
}
