package main

import (
	"fmt"
	"net-cat/pkg"
	"os"
	"sync"
)

var mutex sync.Mutex

func main() {
	input := os.Args
	switch len(input) {
	case 1:
		pkg.ServerCase1()
	case 2:
		pkg.ServerCase2()
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
}
