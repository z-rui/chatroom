package main

import (
	"bufio"
	"github.com/gizak/termui"
	"log"
	"net"
	"os"
	"strings"
	"unicode/utf8"
)

var serverAddr = "127.0.0.1:1357"

var conn net.Conn

const MESSAGE_LOG_LINES = 13

var messageLog []string

func incomingDaemon() {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		recvMessage(line)
	}
}

func recvMessage(msg string) {
	if len(messageLog) == MESSAGE_LOG_LINES {
		messageLog = messageLog[1:]
	}
	messageLog = append(messageLog, msg)
	parMessages.Text = strings.Join(messageLog, "")
	termui.Render(parMessages)
}

func sendMessage(msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

var (
	parMessages *termui.Par
	parEntry    *termui.Par
)

func main() {
	err := termui.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer termui.Close()

	if len(os.Args) == 2 {
		serverAddr = os.Args[1]
	}

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected to", serverAddr)

	createUI()

	go incomingDaemon()

	termui.Loop()
}

func createUI() {
	parMessages = termui.NewPar("")
	parMessages.BorderLabel = "Messages"
	parMessages.Height = 15

	parEntry = termui.NewPar("")
	parEntry.BorderLabel = "Entry"
	parEntry.Height = 9

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, parMessages)),
		termui.NewRow(
			termui.NewCol(12, 0, parEntry)))

	termui.Body.Align()
	termui.Render(termui.Body) // feel free to call Render, it's async and non-block

	termui.Handle("/sys/kbd/<backspace>", func(termui.Event) {
		t := parEntry.Text
		if len(t) == 0 {
			return
		}
		_, size := utf8.DecodeLastRuneInString(t)
		parEntry.Text = t[:len(t)-size]
		termui.Render(parEntry)
	})
	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		sendMessage(parEntry.Text + "\n")
		parEntry.Text = ""
		termui.Render(parEntry)
	})
	termui.Handle("/sys/kbd", func(e termui.Event) {
		t := e.Data.(termui.EvtKbd).KeyStr
		if t == "<space>" {
			t = " "
		} else if utf8.RuneCountInString(t) != 1 {
			return
		}
		parEntry.Text = parEntry.Text + t
		termui.Render(parEntry)
	})
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Render(termui.Body)
	})
}
