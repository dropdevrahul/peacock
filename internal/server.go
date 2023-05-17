// Package gocache this package provides the Server & ServerSettings types
// to start to handle a tcp request
package gocache

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/dropdevrahul/gocache/protocol"
)

const (
	Success = iota + 1
	EmptyValue
	NotFound
	InvalidCommand
	Failure
)

type Server struct {
	c        *Cache
	settings ServerSettings
}

type ServerSettings struct {
	MaxCapacity uint64
}

func NewServer(s ServerSettings) *Server {
	c := NewCache(s.MaxCapacity)

	return &Server{
		c:        c,
		settings: s,
	}
}

type Response struct {
	Data []byte
}

func (s *Server) Serve(conn net.Conn) {
	r := bufio.NewReader(conn)
	h := protocol.Header{}

	err := protocol.ReadHeaders(r, &h)
	if err != nil {
		log.Println(err)
		return
	}

	if h.Len <= 0 {
		return
	}

	rBuff, err := protocol.ReadBody(r, h.Len)
	if err != nil {
		log.Println(err)
		return
	}

	s.Handle(rBuff, conn)
	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) SendResponse(r *Response, conn net.Conn) {
	header := protocol.Header{
		Len: len(r.Data),
	}
	b := append(header.ToBytes(), r.Data...)

	log.Printf("response %s", string(b))

	n, err := conn.Write(b)
	if err != nil || n != len(b) {
		log.Panic("Unable to write complete data to conn", err)
	}
}

func (s *Server) GetPayload(message []byte) []byte {
	p := protocol.ReadPayload(message)
	return p
}

func (s *Server) GetKey(message []byte) string {
	return protocol.GetKey(message)
}

// Get sets response.Data \x00 null byte if key is not present else sets it as value of key
func (s *Server) Get(message []byte, r *Response) {
	key := s.GetKey(message)
	if value, ok := s.c.Get(&key); ok {
		r.Data = []byte(value)
	} else {
		r.Data = []byte("\x00")
	}
}

func (s *Server) Set(message []byte, r *Response) {
	key := s.GetKey(message)
	payload := s.GetPayload(message)
	val := s.c.Set(&key, payload)
	r.Data = []byte(strconv.Itoa(val))
}

func (s *Server) Handle(message []byte, conn net.Conn) {
	r := Response{}
	cmd := protocol.GetCmd(message)
	switch cmd {
	case "SET":
		s.Set(message, &r)
	case "GET":
		s.Get(message, &r)
	case "SET_TTL":
		s.SetTTL(message, &r)
	case "GET_TTL":
		s.GetTTL(message, &r)
	default:
		r.Data = []byte{}
	}

	s.SendResponse(&r, conn)
}

func (s *Server) SetTTL(message []byte, r *Response) {
	key := s.GetKey(message)
	payload := s.GetPayload(message)

	d, err := strconv.Atoi(string(payload))
	if err != nil {
		r.Data = []byte(strconv.Itoa(-1))
		return
	}

	ds := time.Second * time.Duration(d)
	res := s.c.SetTTL(&key, &ds)
	r.Data = []byte(fmt.Sprintf("%d", res))
}

func (s *Server) GetTTL(message []byte, r *Response) {
	key := s.GetKey(message)
	res := s.c.GetTTL(&key)
	ttl := fmt.Sprintf("%d", int(res.Seconds()))
	log.Printf("get ttl %s", ttl)
	r.Data = []byte(ttl)
}
