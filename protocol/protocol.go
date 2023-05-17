package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

const (
	CommandLength    = 11
	KeyLength        = 64
	MaxRequestLength = 2048
	MaxPayloadLength = 2048 - 75
)

type Header struct {
	Len int // denotes the length of payload
}

const HeaderSeparator = "\n"

func (h *Header) ToBytes() []byte {
	return []byte(strconv.Itoa(h.Len) + "\n")
}

func ReadHeaders(r *bufio.Reader, h *Header) error {
	// since tcp is stream based we are never sure when the payload ends/start so requires sending a header with payload length to parse it
	// get payload length from header, header is the first line
	header, err := r.ReadBytes('\n')
	if err != nil {
		return err
	}

	ls := string(header)
	ls = strings.TrimSuffix(ls, "\n")
	l, err := strconv.Atoi(ls)
	if err != nil {
		return err
	}

	h.Len = l
	return nil
}

func ReadBody(r *bufio.Reader, l int) ([]byte, error) {
	rBuff := make([]byte, l)
	n, err := r.Read(rBuff)
	if err != nil {
		return []byte{}, err
	}
	if n != l {
		err = errors.New("error while reading request")
		return []byte{}, err
	}

	return rBuff, nil
}

func ReadPayload(message []byte) []byte {
	return bytes.TrimSpace(message[CommandLength+KeyLength:])
}

func GetKey(message []byte) string {
	key := string(bytes.TrimSpace(message[CommandLength:KeyLength]))
	return key
}

func GetCmd(message []byte) string {
	cmd := string(bytes.ToUpper(bytes.TrimSpace(message[:CommandLength])))
	return cmd
}
