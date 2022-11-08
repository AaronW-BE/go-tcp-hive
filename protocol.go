package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

type Protocol struct {
	length  uint32
	content []byte
}

func Packet(content string) []byte {
	buf := make([]byte, 4+len(content))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(content)))
	copy(buf[4:], content)
	return buf
}

func UnPacket(conn net.Conn) (*Protocol, error) {
	p := &Protocol{}
	header := make([]byte, 4)

	_, err := io.ReadFull(conn, header)
	if err != nil {
		return p, nil
	}

	p.length = binary.BigEndian.Uint32(header)

	content := make([]byte, p.length)
	_, err = io.ReadFull(conn, content)
	if err != nil {
		return p, err
	}
	p.content = content
	return p, nil
}

func (p *Protocol) parseContent() (map[string]interface{}, error) {
	var object map[string]interface{}

	err := json.Unmarshal(p.content, &object)
	if err != nil {
		return object, err
	}
	return object, nil
}
