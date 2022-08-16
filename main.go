package main

import (
	"fmt"
	"net"

	"github.com/dropdevrahul/gocache/gocache"
)

func main() {
	PORT := "0.0.0.0:8888"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		server := &gocache.Server{
			Conn: c,
		}
		go server.Serve()
	}
}
