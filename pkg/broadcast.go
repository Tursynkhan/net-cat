package pkg

import (
	"fmt"
	"sync"
)

func Broadcaster(mutex *sync.Mutex) {
	for {
		select {
		case msg := <-join:
			mutex.Lock()
			for _, client := range clients {
				if client.name == msg.userName {
					for _, w := range historytext {
						fmt.Fprintf(client.conn, "%s", Clear(""))
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
