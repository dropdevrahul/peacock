package gocache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
)

const (
	COMMAND_LENGTH             int = 11 // in bytes
	KEY_LENGTH                 int = 64 // in bytes
	MAX_PAYLOAD_SIZE               = 1468
	HEADER_CONTENT_LENGTH_NAME     = "CONTENT-LENGTH:"
)

const (
	SUCCESS = iota + 1
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
	// since tcp is stream based we are never sure when the payload ends/start so requires sending a header with payload length to parse it
	// also since ordering is not guranteed we need a way to track requests so its better to do it in a channel with timeout
	// we will use fixed length data load in order to gurantee we process one message at a time without overlap
	rBuff := make([]byte, 1468)
	_, _ = bufio.NewReader(s.Conn).Read(rBuff)

	err := s.Handle(rBuff)

	return err
}

func (s *Server) SendResponse(r *Response) {
	header := fmt.Sprintf("%d %s\n", r.Status, r.Error)
	b := append([]byte(header), r.Data...)
	s.Conn.Write(b)
}

func (s *Server) Set(message []byte) (CacheData, error) {
	key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))

	payload := bytes.TrimSpace(message[COMMAND_LENGTH+KEY_LENGTH:])
	err := HashMapCache.Set(&key, payload)
	if err != nil {
		return CacheData{}, err
	}
	val, _ := HashMapCache.Get(&key)
	return val, nil
}

func (s *Server) Handle(message []byte) error {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))
	if cmd == "SET" {
		val, err := s.Set(message)
		if err != nil {
			r := Response{}
			r.Status = FAILURE
			r.Error = "Empty value for key"
			s.SendResponse(&r)
		}

		r := Response{
			Status: SUCCESS,
			Data:   val.BytesData,
		}

		s.SendResponse(&r)

		return nil
	} else if cmd == "GET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))
		if value, ok := HashMapCache.Get(&key); ok {
			r := Response{}
			payload := value.BytesData
			r.Status = SUCCESS
			r.Error = ""
			r.Data = payload
			s.SendResponse(&r)
			return nil
		}

		r := Response{}
		r.Status = FAILURE
		r.Error = "not found"
		s.SendResponse(&r)

		return nil
	}

	return errors.New("Unkown command:" + cmd)
}
