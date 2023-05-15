// Package gocache this package provides the Server & ServerSettings types
// to start to handle a tcp request
package gocache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	CommandLength    = 11
	KeyLength        = 64
	MaxRequestLength = 2048
	MaxPayloadLength = 2048 - 75
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
	Status int
	Error  string
	Data   []byte
}

func (s *Server) Serve(conn net.Conn) {

	reader := bufio.NewReader(conn)

	// since tcp is stream based we are never sure when the payload ends/start so requires sending a header with payload length to parse it
	// get payload length from header, header is the first line
	header, err := reader.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return
	}

	h := string(header)
	h = strings.TrimSuffix(h, "\n")
	l, err := strconv.Atoi(h)
	if err != nil {
		log.Println(err)
		return
	}

	rBuff := make([]byte, l)
	n, err := reader.Read(rBuff)
	if err != nil {
		log.Println(err)
		return
	}

	if n != l {
		err = errors.New("error while reading request")
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
	header := fmt.Sprintf("%d %s\n", r.Status, r.Error)
	b := append([]byte(header), r.Data...)
	n, err := conn.Write(b)
	if err != nil || n != len(b) {
		log.Panic("Unable to write complete data to conn", err)
	}
}

func (s *Server) GetKey(message []byte) string {
	key := string(bytes.TrimSpace(message[CommandLength:KeyLength]))

	return key
}

func (s *Server) Get(message []byte, r *Response) {
	key := s.GetKey(message)
	if value, ok := s.c.Get(&key); ok {
		r.Status = Success
		r.Error = ""
		r.Data = []byte(value)
	} else {
		r.Status = NotFound
		r.Error = "not found"
	}
}

func (s *Server) Set(message []byte, r *Response) {
	key := s.GetKey(message)
	payload := bytes.TrimSpace(message[CommandLength+KeyLength:])
	err := s.c.Set(&key, payload)
	if err != nil {
		fmt.Println(err)
		r.Status = EmptyValue // TODO handle more errors
		r.Error = err.Error()
		r.Data = []byte{}
	} else {
		val, _ := s.c.Get(&key)
		r.Status = Success
		r.Data = []byte(val)
		r.Error = ""
	}
}

func (s *Server) Handle(message []byte, conn net.Conn) {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:CommandLength])))
	r := Response{}
	switch cmd {
	case "SET":
		s.Set(message, &r)
	case "GET":
		s.Get(message, &r)
	case "SET_TTL":
		s.SetTTL(message, &r)
	default:
		r.Error = "invalid command"
		r.Status = InvalidCommand
	}

	s.SendResponse(&r, conn)
}

func (s *Server) SetTTL(message []byte, r *Response) {
	key := s.GetKey(message)
	payload := bytes.TrimSpace(message[CommandLength+KeyLength:])

	d, err := strconv.Atoi(string(payload))
	if err != nil {
		r.Data = []byte{}
		r.Status = InvalidCommand
		r.Error = fmt.Sprintf("invalid ttl value %s", string(payload))

		return
	}

	ds := time.Second * time.Duration(d)
	res := s.c.SetTTL(&key, &ds)

	r.Data = []byte(fmt.Sprintf("%d", res))
	r.Status = Success
}
