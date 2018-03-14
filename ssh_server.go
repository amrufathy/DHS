package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
	"strings"
	"os/exec"
	"sync"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		command := strings.Join(s.Command()," ")
		fmt.Println(command)
		parts := strings.Fields(command)
		out, err := exec.Command(parts[0], parts[1:]...).Output()
		if err != nil {
			fmt.Println("error occured")
			fmt.Printf("%s\n", err)
		}
		fmt.Printf("%s\n", out)
		fmt.Println(command)
		io.WriteString(s, fmt.Sprintf(string(out)))
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0],parts[1]).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	//wg.Done()
}