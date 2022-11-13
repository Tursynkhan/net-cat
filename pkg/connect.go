package pkg

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

func HandleConnection(conn net.Conn, mutex *sync.Mutex) {
	Welcome(conn)
	username := GetName(conn)
	tempClient := Client{
		name: username,
		addr: conn.RemoteAddr().String(),
		conn: conn,
	}
	mutex.Lock()
	clients[username] = tempClient
	if len(clients) > 10 {
		tempClient.conn.Write([]byte("Chat is full!"))
		delete(clients, username)
		tempClient.conn.Close()
		mutex.Unlock()
		return
	}

	mutex.Unlock()
	t := time.Now().Format("2006-01-02 15:04:05")
	mutex.Lock()
	join <- NewMessage("has joined our chat...", conn, tempClient, t)
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
		if IsCorrect(input.Text(), conn, time, username) == false {
			continue
		}
		text := fmt.Sprintf("[%s][%s]:%s\n", time, username, input.Text())
		mutex.Lock()
		historytext = append(historytext, text)
		mutex.Unlock()
		mutex.Lock()
		messages <- NewMessage(input.Text(), conn, tempClient, time)
		mutex.Unlock()
	}
	mutex.Lock()
	delete(clients, username)
	leaving <- NewMessage("has left our chat...", conn, tempClient, t)
	conn.Close()
	mutex.Unlock()
}
