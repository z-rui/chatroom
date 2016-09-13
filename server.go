package main

import (
	"bufio"
	"container/list"
	"io"
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

func serve(conn net.Conn) {
	clients.Lock()
	e := clients.PushBack(conn)
	clients.Unlock()
	log.Println(conn.RemoteAddr(), "connected")

	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			if err == io.EOF {
				break
			} else if err, ok := err.(net.Error); ok && !err.Temporary() {
				break
			}
			continue
		}
		broadcast(line)
	}

	clients.Lock()
	clients.Remove(e)
	clients.Unlock()
	log.Println(conn.RemoteAddr(), "disconnected")
}

func broadcast(line []byte) {
	clients.Lock()
	defer clients.Unlock()
	for e := clients.Front(); e != nil; e = e.Next() {
		go func(conn net.Conn) {
			_, err := conn.Write(line)
			if err != nil {
				log.Println(err.Error())
			}
		}(e.Value.(net.Conn))
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
		go serve(conn)
	}
}
