package gocache

import (
	"bufio"
	"bytes"
	"errors"
	"net"
)

const COMMAND_LENGTH int = 11 // in bytes
const KEY_LENGTH int = 64     // in bytes
const MAX_PAYLOAD_SIZE = 1468
const HEADER_CONTENT_LENGTH_NAME = "CONTENT-LENGTH:"

type Server struct {
	Conn net.Conn
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

func (s *Server) SendResponse(m []byte) {
	s.Conn.Write(m)
}

func (s *Server) Handle(message []byte) error {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))
	if cmd == "SET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))

		payload := bytes.TrimSpace(message[COMMAND_LENGTH+KEY_LENGTH:])

		HashMapCache.Set(&key, payload)
		val, _ := HashMapCache.Get(&key)

		s.SendResponse(val.BytesData)
	} else if cmd == "GET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))

		if value, ok := HashMapCache.Get(&key); ok {
			payload := value.BytesData
			s.Conn.Write(payload)
		}

		s.Conn.Write([]byte{})
	}

	return errors.New("Unkown command:" + cmd)
}
