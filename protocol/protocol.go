package protocol

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
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

func ReadPayload(r *bufio.Reader, l int) ([]byte, error) {
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
