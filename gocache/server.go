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

func (s *Server) Handle(message []byte) {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))

	if cmd == "SET" {
		val, err := s.Set(message)
		if err != nil {
			r := Response{}
			r.Status = EMPTYVALUE
			r.Error = "Empty value for key"
			s.SendResponse(&r)

			return
		}

		r := Response{
			Status: SUCCESS,
			Data:   val.BytesData,
		}

		s.SendResponse(&r)

		return

	} else if cmd == "GET" {
		key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))
		if value, ok := HashMapCache.Get(&key); ok {
			r := Response{}
			payload := value.BytesData
			r.Status = SUCCESS
			r.Error = ""
			r.Data = payload
			s.SendResponse(&r)
			return
		}

		r := Response{}
		r.Status = NOTFOUND
		r.Error = "not found"
		s.SendResponse(&r)

		return
	}

	r := Response{
		Status: INVALIDCOMMAND,
		Error:  "invalid command",
	}

	s.SendResponse(&r)
}
