package main

import (
	"net"
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"encoding/binary"
	"strconv"
	"math/rand"
)

var (
	sAddr       string
	cNumReaders int
	cNumWriters int
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
	defer conn.Close()

	for i := 0; i < cNumReaders; i++ {
		go readerClient(i)
	}

	for i := 0; i < cNumWriters; i++ {
		go writerClient(i)
	}

	for {
	}
}

func readerClient(idx int) {
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	var number int64

	for {
		// request
		conn.Write([]byte("read\n"))

		err := binary.Read(conn, binary.BigEndian, &number)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Reader #", idx, "=> value from server", number)
	}
}

func writerClient(idx int) {
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	for {
		numToWrite := rand.Intn(100)
		conn.Write([]byte("write " + strconv.Itoa(numToWrite) + "\n"))

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Writer #", idx, "writing", numToWrite)
		fmt.Println("Message from server", message)
	}
}
