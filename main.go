package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	CONN_PORT = ":8989"
	CONN_TYPE = "tcp"
)

func main() {
	input := os.Args
	switch len(input) {
	case 1:
		l, err := net.Listen(CONN_TYPE, CONN_PORT)
		if err != nil {
			log.Fatalf("unable to start server: %s", err.Error())
		}

		defer l.Close()
		log.Printf("Listening on the port :8989")
		go broadcaster(&mutex)
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("unable to accept connection: %s", err.Error())
				conn.Close()
				continue
			}
			go handleConnection(conn, &mutex)
		}
	case 2:

		arg := os.Args[1]
		port := fmt.Sprintf(":%s", arg)
		l, err := net.Listen(CONN_TYPE, port)
		if err != nil {
			log.Fatalf("unable to start server: %s", err.Error())
		}

		defer l.Close()
		fmt.Printf("Listening on the port %s\n", port)
		go broadcaster(&mutex)
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("unable to accept connection: %s", err.Error())
				conn.Close()
				continue
			}
			go handleConnection(conn, &mutex)
		}
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
}

var (
	mutex       sync.Mutex
	clients     = make(map[string]Client)
	historytext = []string{}
	leaving     = make(chan Message)
	messages    = make(chan Message)
	join        = make(chan Message)
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
	history  []string
}

func handleConnection(conn net.Conn, mutex *sync.Mutex) {
	welcome(conn)
	username := getName(conn)

	tempClient := Client{
		name: username,
		addr: conn.RemoteAddr().String(),
		conn: conn,
	}
	mutex.Lock()
	clients[username] = tempClient
	if len(clients) > 10 {
		tempClient.conn.Write([]byte("Maximum 10 users"))
		delete(clients, username)
		tempClient.conn.Close()
		mutex.Unlock()
		return
	}

	mutex.Unlock()
	t := time.Now().Format("2006-01-02 15:04:05")
	mutex.Lock()
	join <- newMessage("has joined our chat...", conn, tempClient, t)
	mutex.Unlock()
	txt := fmt.Sprintf("[%s][%s]:", t, username)
	fmt.Fprintf(conn, txt)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		time := time.Now().Format("2006-01-02 15:04:05")
		if input.Text() == "" {
			fmt.Fprintf(conn, "you can't send empty messages\n")
			fmt.Fprintf(conn, "[%s][%s]:", time, username)
			continue
		}
		text := fmt.Sprintf("[%s][%s]:%s\n", time, username, input.Text())
		mutex.Lock()
		historytext = append(historytext, text)
		mutex.Unlock()
		mutex.Lock()
		messages <- newMessage(input.Text(), conn, tempClient, time)
		mutex.Unlock()
	}
	mutex.Lock()
	delete(clients, username)
	leaving <- newMessage("has left our chat...", conn, tempClient, t)
	conn.Close()
	mutex.Unlock()
}

func newMessage(msg string, conn net.Conn, cl Client, time string) Message {
	return Message{
		text:     msg,
		address:  cl.addr,
		userName: cl.name,
		time:     time,
		history:  historytext,
	}
}

func broadcaster(mutex *sync.Mutex) {
	for {
		select {
		case msg := <-join:
			mutex.Lock()
			for _, client := range clients {
				if client.name == msg.userName {
					for _, w := range historytext {
						fmt.Fprintf(client.conn, "%s", clear(""))
						fmt.Fprintf(client.conn, "%s[%s][%s]:", w, msg.time, client.name)
					}
				}
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n%s %s\n[%s][%s]:", msg.userName, msg.text, msg.time, client.name)
				}
			}
			mutex.Unlock()
		case msg := <-messages:
			mutex.Lock()
			for _, client := range clients {
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n[%s][%s]:%s\n", msg.time, msg.userName, msg.text)
				}
				fmt.Fprintf(client.conn, "[%s][%s]:", msg.time, client.name)
			}
			mutex.Unlock()
		case msg := <-leaving:
			mutex.Lock()
			for _, client := range clients {
				if client.name != msg.userName {
					fmt.Fprintf(client.conn, "\n%s %s\n[%s][%s]:", msg.userName, msg.text, msg.time, client.name)
				}
			}
			mutex.Unlock()
		}
	}
}

func welcome(conn net.Conn) {
	file, err := os.ReadFile("penguin.txt")
	if err != nil {
		fmt.Printf("couldn't read this file")
	}
	strWelcome := (string(file))
	conn.Write([]byte("Welcome to TCP-Chat!\n" + strWelcome + "\n"))
}

func getName(conn net.Conn) string {
	conn.Write([]byte("[Enter your name]:"))
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}

	temp := strings.TrimSpace(string(data))
	if temp == "" || len(temp) == 0 {
		return getName(conn)
	}
	for i, _ := range clients {
		if i == temp {
			fmt.Fprintf(conn, "User already exist\n")
			return getName(conn)
		}
	}
	return temp
}

func clear(a string) string {
	return "\r" + strings.Repeat(" ", len(a)) + "\r"
}
