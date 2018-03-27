package main

import (
	"net"
	"fmt"
	"github.com/go-ini/ini"
	"math/rand"
	"log"
	"os"
	"sync"
	"net/rpc"
)

var (
	globalNumber = rand.Int63n(100)
	addr         string
	READER       *log.Logger
	WRITER       *log.Logger
	sSequence    int64 = 0
	sNumReaders        = 0
	sNumWriters        = 0
	mutex        sync.Mutex // mutex for R/W data
	mSeq         sync.Mutex // mutex for sequence number
	mReaders     sync.Mutex // mutex for number of readers
	mWriters     sync.Mutex // mutex for number of writers
)

type Data int

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
	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()
	fmt.Println("Server is listening on", addr)

	data := new(Data)
	rpc.Register(data)

	for {
		// accept a connection
		conn, _ := listener.AcceptTCP()

		// Print client's remote address
		remoteAddr := conn.RemoteAddr().String()
		fmt.Println("Client connected from", remoteAddr)

		if err != nil {
			fmt.Printf("Connection error %s\n", err)
		}

		rpc.ServeConn(conn)
	}
}

func (d *Data) Read(rIdx int, reply *ReadStruct) error {
	increaseNumReaders()
	defer decreaseNumReaders()
	defer increaseSequenceNumber()

	//time.Sleep(time.Millisecond * time.Duration(500))

	reply.Result = globalNumber
	reply.Rseq = sSequence

	// log
	READER.Printf("%d\t%d\t%d\t%d\n", sSequence, globalNumber, sNumReaders, rIdx)

	return nil
}

func (d *Data) Write(ws WriteStruct, wSeq *int64) error {
	increaseNumWriters()
	defer decreaseNumWriters()
	defer increaseSequenceNumber()

	//time.Sleep(time.Millisecond * time.Duration(500))

	// log
	WRITER.Printf("%d\t%d\t%d\t%d\n", sSequence, globalNumber, ws.NewVal, ws.Widx)

	*wSeq = sSequence

	mutex.Lock()
	globalNumber = ws.NewVal
	mutex.Unlock()

	return nil
}

func increaseSequenceNumber() {
	mSeq.Lock()
	sSequence++
	mSeq.Unlock()
}

func increaseNumReaders() {
	mReaders.Lock()
	sNumReaders++
	mReaders.Unlock()
}

func decreaseNumReaders() {
	mReaders.Lock()
	sNumReaders--
	mReaders.Unlock()
}

func increaseNumWriters() {
	mWriters.Lock()
	sNumWriters++
	mWriters.Unlock()
}

func decreaseNumWriters() {
	mWriters.Lock()
	sNumWriters--
	mWriters.Unlock()
}
