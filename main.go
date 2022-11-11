package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_PORT = ":8989"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer l.Close()
	log.Printf("started server on :8989")

	go broadcaster()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			conn.Close()
			continue
		}
		go handleConnection(conn)

	}
}

var (
	clients  = make(map[string]Client)
	leaving  = make(chan Message)
	messages = make(chan Message)
	join     = make(chan Message)
)

type Client struct {
	name string
	addr string
	conn net.Conn
}

type Message struct {
	text     string
	address  string
	userName string
	time     string
}

func handleConnection(conn net.Conn) {
	welcome(conn)
	username := getName(conn)
	tempClient := Client{
		name: username,
		addr: conn.RemoteAddr().Network(),
		conn: conn,
	}
	clients[username] = tempClient
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(conn, "[%s][%s]:", t, username)
	join <- newMessage("has joined our chat...", conn, tempClient, t)

	input := bufio.NewScanner(conn)
	for input.Scan() {

		time := time.Now().Format("2006-01-02 15:04:05")
		messages <- newMessage(input.Text(), conn, tempClient, time)
	}
	delete(clients, conn.RemoteAddr().String())

	leaving <- newMessage("has left our chat...", conn, tempClient, t)

	conn.Close()
}

func newMessage(msg string, conn net.Conn, cl Client, time string) Message {
	return Message{
		text:     msg,
		address:  cl.addr,
		userName: cl.name,
		time:     time,
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-join:
			for _, client := range clients {
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n%s %s\n[%s][%s]:", msg.userName, msg.text, msg.time, client.name)
				}
			}
		case msg := <-messages:
			for _, client := range clients {
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n[%s][%s]: %s\n", msg.time, msg.userName, msg.text)
				}
				fmt.Fprintf(client.conn, "[%s][%s]: ", msg.time, client.name)
			}
		case msg := <-leaving:
			for _, client := range clients {
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n%s %s\n[%s][%s]:", msg.userName, msg.text, msg.time, client.name)
				}
			}
		}
	}
}

func welcome(conn net.Conn) {
	file, err := os.ReadFile("text.txt")
	if err != nil {
		fmt.Printf("couldn't read this file")
	}
	strWelcome := (string(file))
	conn.Write([]byte("Welcome to TCP-Chat!\n" + strWelcome + "\n"))
	conn.Write([]byte("[Enter your name]:"))
}

func getName(conn net.Conn) string {
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}

	temp := strings.TrimSpace(string(data))
	return temp
}

func clear(a string) string {
	return "\r" + strings.Repeat(" ", len(a)) + "\r"
}
