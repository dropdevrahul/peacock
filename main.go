// Entry point for the server
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
	size := flag.Uint64("max-size", 100000,
		"Maximum number of items allowed in cache before LRU kicks in")

	flag.Parse()

	l, err := net.Listen("tcp4", *host+":"+*port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	server := gocache.NewServer(
		gocache.ServerSettings{
			MaxCapacity: *size,
		},
	)

	fmt.Printf("Server listening on %s:%s\n", *host, *port)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go server.Serve(c)
	}
}
