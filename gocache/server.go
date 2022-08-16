package gocache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const COMMAND_LENGTH int = 11            // in bytes
const KEY_LENGTH int = 64                // in bytes
const MAX_PAYLOAD_LENGTH int = 10 * 1024 // max 10 KB
const NULL_BYTE byte = 0
const TOTAL_PAYLOAD_BYTES int = COMMAND_LENGTH + KEY_LENGTH + MAX_PAYLOAD_LENGTH + 1
const MIN_BODY_SIZE int = COMMAND_LENGTH + KEY_LENGTH + 1
const HEADER_CONTENT_LENGTH_NAME = "CONTENT-LENGTH:"

func clen(n []byte) int {
	var i int = 0
	var nlen = len(n)
	for ; i < nlen; i++ {
		if n[i] == 0 {
			return i // we don't want i + 1 since we want to skip null byte
		}
	}
	if nlen <= MAX_PAYLOAD_LENGTH {
		return nlen
	}
	return MAX_PAYLOAD_LENGTH
}

type Server struct {
	Conn net.Conn
}

func (s *Server) Serve() error {
	defer s.Conn.Close()
	// Try to get the content length
	// we will assume the content-length is contained in first line of the message
	msg, _ := bufio.NewReader(s.Conn).ReadBytes('\n')
	contLen, err := s.ParseHeader(msg)
	if err != nil {
		return err
	}
	fmt.Println("len", contLen)
	s.Conn.Write([]byte("ACK"))

	buff := make([]byte, contLen)
	_, err = bufio.NewReader(s.Conn).Read(buff)

	// valid content length, read buffer til content-length
	//msg, err = bufio.NewReader(s.Conn).ReadBytes(NULL_BYTE)
	if err != nil {
		panic(err)
		return err
	}
	err = s.Handle(buff)
	return err
}

func (s *Server) ParseHeader(message []byte) (int, error) {
	header := string(bytes.ToUpper(bytes.TrimSpace(message)))
	splits := strings.Split(header, HEADER_CONTENT_LENGTH_NAME)
	if len(splits) != 2 {
		return -1, errors.New("Invalid header")
	}
	if splits[0] != "" {
		return -1, errors.New("Invalid header")
	}
	contentLen, err := strconv.Atoi(splits[1])
	if contentLen < MIN_BODY_SIZE {
		return contentLen, errors.New("Invalid content length")
	}
	return contentLen, err
}

func (s *Server) Handle(message []byte) error {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))
	if cmd == "SET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))

		payload := message[COMMAND_LENGTH+KEY_LENGTH:]
		last_index := clen(payload)

		payload = payload[:last_index]
		HashMapCache.Set(&key, payload)
		val, _ := HashMapCache.Get(&key)
		fmt.Println(string(val.bytesData))
		s.Conn.Write([]byte("OK"))
	} else if cmd == "GET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))
		if value, ok := HashMapCache.Get(&key); ok {
			last_index := clen(value.bytesData)
			payload := value.bytesData[:last_index]
			s.Conn.Write(payload)
		}
		fmt.Println("No key found")
	}
	return errors.New("Unkown command:" + cmd)
}
