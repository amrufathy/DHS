package main

import (
	"net"
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"encoding/binary"
	"strconv"
	"math/rand"
	"sync"
	"os"
	"log"
)

var (
	sAddr       string
	cNumReaders int
	cNumWriters int
	cNumAccess  int
	wg          sync.WaitGroup
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
	cAccess, _ := cInfo.GetKey("num_accesses")

	cReadersInt, _ := strconv.Atoi(cReaders.String())
	cWritersInt, _ := strconv.Atoi(cWriters.String())
	cAccessInt, _ := strconv.Atoi(cAccess.String())

	cNumReaders = cReadersInt
	cNumWriters = cWritersInt
	cNumAccess = cAccessInt
}

func main() {
	fmt.Println("Requesting a connection from server...")

	// send request to server for connection
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	for i := 0; i < cNumReaders; i++ {
		wg.Add(1)
		go readerClient(i)
	}

	for i := 0; i < cNumWriters; i++ {
		wg.Add(1)
		go writerClient(i)
	}

	wg.Wait()
}

func readerClient(idx int) {
	defer wg.Done()
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	logFile, _ := os.Create("./logs/readers/reader" + strconv.Itoa(idx) + ".log")
	logger := log.New(logFile, "", log.Ltime)

	var number int64

	for i := 0; i < cNumAccess; i++ {
		// request
		conn.Write([]byte("read\n"))

		err := binary.Read(conn, binary.BigEndian, &number)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Reader #", idx, "=> value from server", number)
		logger.Printf("%d\t%d\n", idx, number)
	}
}

func writerClient(idx int) {
	defer wg.Done()
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	logFile, _ := os.Create("./logs/writers/writer" + strconv.Itoa(idx) + ".log")
	logger := log.New(logFile, "", log.Ltime)

	for i := 0; i < cNumAccess; i++ {
		numToWrite := rand.Intn(100)
		conn.Write([]byte("write " + strconv.Itoa(numToWrite) + "\n"))

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Writer #", idx, "writing", numToWrite)
		fmt.Println("Message from server", message)

		logger.Printf("%d\t%d\n", idx, numToWrite)
	}
}
