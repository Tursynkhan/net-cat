package pkg

import (
	"net"
	"sync"
)

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
