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
)

const (
	COMMAND_LENGTH     int = 11 // in bytes
	KEY_LENGTH         int = 64 // in bytes
	MAX_REQUEST_LENGTH     = 2048
	MAX_PAYLOAD_LENGTH     = 2048 - 75
)

const (
	SUCCESS = iota + 1
	EMPTYVALUE
	NOTFOUND
	INVALIDCOMMAND
	FAILURE
)

type Server struct {
	Conn net.Conn
}

type Response struct {
	Status int
	Error  string
	Data   []byte
}

func (s *Server) Serve() error {
	defer s.Conn.Close()
	reader := bufio.NewReader(s.Conn)

	// since tcp is stream based we are never sure when the payload ends/start so requires sending a header with payload length to parse it
	// get payload length from header, header is the first line
	header, err := reader.ReadBytes('\n')
	if err != nil {
		log.Println(err)

		return errors.New("Error while reading request")
	}

	h := string(header)
	h = strings.TrimSuffix(h, "\n")
	l, err := strconv.Atoi(h)
	if err != nil {
		log.Println(err)

		return errors.New("Error reading payload length")
	}

	rBuff := make([]byte, l)
	n, err := reader.Read(rBuff)
	if err != nil {
		log.Println(err)

		return errors.New("Error while reading request")
	}

	if n != l {
		err = errors.New("Error while reading request")
		log.Println(err)

		return err
	}

	s.Handle(rBuff)

	return nil
}

func (s *Server) SendResponse(r *Response) {
	header := fmt.Sprintf("%d %s\n", r.Status, r.Error)
	b := append([]byte(header), r.Data...)
	s.Conn.Write(b)
}

func (s *Server) GetKey(message []byte) string {
	key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))

	return key
}

func (s *Server) Del(message []byte, r *Response) {
	key := s.GetKey(message)
	if value, ok := HashMapCache.Get(&key); ok {
		HashMapCache.Del(&key)
		r.Status = SUCCESS
		r.Error = ""
		r.Data = value.BytesData
	} else {
		r.Status = NOTFOUND
		r.Error = "not found"
	}
}

func (s *Server) Get(message []byte, r *Response) {
	key := s.GetKey(message)

	if value, ok := HashMapCache.Get(&key); ok {
		r.Status = SUCCESS
		r.Error = ""
		r.Data = value.BytesData
	} else {
		r.Status = NOTFOUND
		r.Error = "not found"
	}
}

func (s *Server) Set(message []byte, r *Response) {
	key := s.GetKey(message)
	payload := bytes.TrimSpace(message[COMMAND_LENGTH+KEY_LENGTH:])

	err := HashMapCache.Set(&key, payload)
	if err != nil {
		r.Status = EMPTYVALUE // TODO handle more errors
		r.Error = err.Error()
		r.Data = []byte{}
	} else {
		val, _ := HashMapCache.Get(&key)
		r.Status = SUCCESS
		r.Data = val.BytesData
		r.Error = ""
	}
}

func (s *Server) Handle(message []byte) {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))
	r := Response{}
	switch cmd {
	case "SET":
		s.Set(message, &r)
	case "GET":
		s.Get(message, &r)
	case "DEL":
		s.Del(message, &r)
	default:
		r.Error = "invalid command"
		r.Status = INVALIDCOMMAND
	}

	s.SendResponse(&r)
}
