package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

const (
	serverAddr = "127.0.0.1:1357"
)

/*
func incomingDaemon(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			continue
		}
		os.Stdout.Write(line)
	}
}
*/

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected to", serverAddr)
	//go incomingDaemon(conn)
	go io.Copy(os.Stdout, conn)
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			log.Fatalln(err.Error())
		}
		conn.Write(line)
	}
}
