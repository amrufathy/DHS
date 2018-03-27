package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"strconv"
	"math/rand"
	"sync"
	"os"
	"log"
	"net/rpc"
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
	conn, _ := rpc.Dial("tcp", sAddr)
	defer conn.Close()

	logFile, _ := os.Create("./logs/readers/reader" + strconv.Itoa(idx) + ".log")
	logger := log.New(logFile, "", log.Ltime)

	for i := 0; i < cNumAccess; i++ {
		rv := new(ReadStruct)

		// read value and send idx
		if err := conn.Call("Data.Read", idx, &rv); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Reader #", idx, "Access #", i, "=> value from server", rv.Result)
		logger.Printf("%d\t%d\t%d\n", idx, rv.Rseq, rv.Result)
	}
}

func writerClient(idx int) {
	defer wg.Done()
	conn, _ := rpc.Dial("tcp", sAddr)
	defer conn.Close()

	logFile, _ := os.Create("./logs/writers/writer" + strconv.Itoa(idx) + ".log")
	logger := log.New(logFile, "", log.Ltime)

	var wSeq int64

	for i := 0; i < cNumAccess; i++ {
		numToWrite := rand.Int63n(100)
		ws := WriteStruct{NewVal: numToWrite, Widx: idx}

		// send new value and get wSeq
		if err := conn.Call("Data.Write", ws, &wSeq); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Writer #", idx, "Access #", i, "=> writing", numToWrite)

		logger.Printf("%d\t%d\t%d\n", idx, numToWrite, wSeq)
	}
}
