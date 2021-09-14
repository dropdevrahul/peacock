package gocache

import (
    "net"
    "bufio"
    "bytes"
    "fmt"
)

const COMMAND_LENGTH int = 11 // in bytes
const KEY_LENGTH int= 64 // in bytes
const MAX_PAYLOAD_LENGTH int= 10*1024 // max 10 KB
const NULL_BYTE byte = 0

func clen(n []byte) int {
    var i int = 0
    var nlen = len(n)
    for ; i < nlen; i++ {
        if n[i] == 0 {
            return i // we don;t want i + 1 since we want to skip null byte 
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


func (s* Server) Read() error {
    for {
        msg, _ := bufio.NewReader(s.Conn).ReadBytes('\n')
        s.Handle(msg)
    }
}


func (s *Server) Handle(message []byte) {
    cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:COMMAND_LENGTH])))
    if cmd == "SET" {
        key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))
        payload := message[COMMAND_LENGTH+KEY_LENGTH:]
        last_index := clen(payload)
        payload = payload[:last_index]
        HashMapCache.Set(&key, payload)
    }
    if cmd == "GET" {
        key := string(bytes.TrimSpace(message[COMMAND_LENGTH:KEY_LENGTH]))
        if value, ok := HashMapCache.Get(&key); ok {
            last_index := clen(value.bytesData)
            payload := value.bytesData[:last_index]
            s.Conn.Write(payload)
        }
        fmt.Println("No key found")
    }
}
