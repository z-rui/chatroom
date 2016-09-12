package main

import (
	"bufio"
	"container/list"
	"log"
	"net"
	"os"
	"sync"
)

var listenAddr = ":1357"

var clients struct {
	list.List
	sync.Mutex
}

func incomingDaemon(e *list.Element) {
	conn := e.Value.(net.Conn)
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			if err, ok := err.(net.Error); ok && !err.Temporary() {
				clients.Lock()
				clients.Remove(e)
				clients.Unlock()
				log.Println(conn.RemoteAddr(), "disconnected")
				return
			}
			continue
		}
		broadcast(line)
	}
}

func broadcast(line []byte) {
	clients.Lock()
	defer clients.Unlock()
	for e := clients.Front(); e != nil; e = e.Next() {
		conn := e.Value.(net.Conn)
		_, err := conn.Write(line)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func main() {
	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Listening", listenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		clients.Lock()
		e := clients.PushBack(conn)
		clients.Unlock()
		go incomingDaemon(e)
		log.Println("Client", conn.RemoteAddr(), "connected")
	}
}