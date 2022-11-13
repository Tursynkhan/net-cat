package pkg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func NewMessage(msg string, conn net.Conn, cl Client, time string) Message {
	return Message{
		text:     msg,
		address:  cl.addr,
		userName: cl.name,
		time:     time,
		history:  historytext,
	}
}

func Welcome(conn net.Conn) {
	file, err := os.ReadFile("penguin.txt")
	if err != nil {
		fmt.Printf("couldn't read this file")
	}
	strWelcome := (string(file))
	conn.Write([]byte("Welcome to TCP-Chat!\n" + strWelcome + "\n"))
}

func GetName(conn net.Conn) string {
	conn.Write([]byte("[Enter your name]:"))
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}

	temp := strings.TrimSpace(string(data))
	if temp == "" || len(temp) == 0 {
		fmt.Fprintln(conn, "Incorrect input")
		return GetName(conn)
	}
	for _, w := range temp {
		if w < 32 || w > 127 {
			fmt.Fprintln(conn, "Incorrect input")
			return GetName(conn)
		}
	}
	for i, _ := range clients {
		if i == temp {
			fmt.Fprintf(conn, "User already exist\n")
			return GetName(conn)
		}
	}
	return temp
}

func Clear(a string) string {
	return "\r" + strings.Repeat(" ", len(a)) + "\r"
}

func IsCorrect(s string, conn net.Conn, time string, username string) bool {
	for _, w := range s {
		if w < 32 || w > 127 {
			fmt.Fprintln(conn, "Incorrect input")
			fmt.Fprintf(conn, "[%s][%s]:", time, username)
			return false
		}
	}
	return true
}
