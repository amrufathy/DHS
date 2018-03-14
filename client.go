package main

import (
	"net"
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"strconv"
	"math/rand"
	"sync"
	"os"
	"log"
	"strings"
	"time"
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
	// send request to server for connection
	//conn, _ := net.Dial("tcp", sAddr)
	//defer conn.Close()

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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var (
		number int64
		rSeq   int64
	)

	for i := 0; i < cNumAccess; i++ {
		// request
		writer.WriteString("read\n")
		writer.Flush()

		// read value
		if line, err := reader.ReadString('\n'); err != nil {
			println("Error: ", err.Error())
		} else {
			val, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
			number = val
		}

		// read sequence number
		if line, err := reader.ReadString('\n'); err != nil {
			println("Error: ", err.Error())
		} else {
			val, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
			rSeq = val
		}

		writer.WriteString(strconv.FormatInt(int64(idx), 10) + "\n")
		writer.Flush()

		fmt.Println("Reader #", idx, "Access #", i, "=> value from server", number)
		logger.Printf("%d\t&d\t%d\n", idx, rSeq, number)
	}
}

func writerClient(idx int) {
	defer wg.Done()
	conn, _ := net.Dial("tcp", sAddr)
	defer conn.Close()

	logFile, _ := os.Create("./logs/writers/writer" + strconv.Itoa(idx) + ".log")
	logger := log.New(logFile, "", log.Ltime)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var wSeq int64

	for i := 0; i < cNumAccess; i++ {
	numToWrite := rand.Int63n(100)

	// request
	writer.WriteString("write \n")
	writer.Flush()


	// FIXME: why this semi-works ?
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))

	// send value
	writer.WriteString(strconv.FormatInt(numToWrite, 10) + "\n")
	writer.Flush()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))

	fmt.Println("Writer #", idx, "Access #", 0, "writing", numToWrite)

	// read sequence number
	if line, err := reader.ReadString('\n'); err != nil {
		println("Error: ", err.Error())
	} else {
		val, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		wSeq = val
		fmt.Println("Message from server", val)
	}

	logger.Printf("%d\t%d\t%d\n", idx, numToWrite, wSeq)
	}
}
