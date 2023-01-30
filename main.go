package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/dropdevrahul/gocache/gocache"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host to listen on")
	port := flag.String("port", "9999", "port to listen on")

	flag.Parse()

	l, err := net.Listen("tcp4", *host+":"+*port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Printf("Server listening on %s:%s\n", *host, *port)
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
