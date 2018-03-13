package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func main() {
	user := "ghannam"
	host := "localhost:2222"
	command := "go run server.go"
	client, session, err := connectToHost(user, host)
	if err != nil {
		panic(err)
	}
	out, err := session.CombinedOutput(command)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	client.Close()
}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	//var pass string
	//fmt.Print("Password: ")
	//fmt.Scanf("%s\n", &pass)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password("")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
