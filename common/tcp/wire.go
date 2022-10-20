package tcp

import (
	"bufio"
	"encoding/json"
	"net"
)

type Wire struct {
	conn net.Conn
	r    *bufio.Reader
}

func NewWire(addr string) (*Wire, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Wire{
		conn: conn,
		r:    bufio.NewReader(conn),
	}, nil
}

func (w *Wire) Write(o interface{}) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = w.conn.Write(data)
	return err
}

func (w *Wire) Read(o interface{}) error {
	data, err := w.ReadBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

func (w *Wire) ReadBytes() ([]byte, error) {
	data, err := w.r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return data[:len(data)-1], nil
}

func (w *Wire) Close() error {
	return w.conn.Close()
}
