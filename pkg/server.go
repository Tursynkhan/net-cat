package pkg

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	CONN_PORT = ":8989"
	CONN_TYPE = "tcp"
)

func ServerCase1() {
	l, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer l.Close()
	log.Printf("Listening on the port :8989")
	go Broadcaster(&mutex)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			conn.Close()
			continue
		}
		go HandleConnection(conn, &mutex)
	}
}

func ServerCase2() {
	arg := os.Args[1]
	port := fmt.Sprintf(":%s", arg)
	l, err := net.Listen(CONN_TYPE, port)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer l.Close()
	fmt.Printf("Listening on the port %s\n", port)
	go Broadcaster(&mutex)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			conn.Close()
			continue
		}
		go HandleConnection(conn, &mutex)
	}
}
