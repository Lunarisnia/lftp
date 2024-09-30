package dsu

import "bufio"

type ServerMemo struct {
	BufferMap   map[string]*bufio.Writer
	TotalLength int
}

func NewServerMemo() *ServerMemo {
	return &ServerMemo{
		BufferMap: make(map[string]*bufio.Writer),
	}
}

type ClientMemo struct {
	BufferMap   map[string]*bufio.Reader
	TotalLength int
}

func NewClientMemo() *ClientMemo {
	return &ClientMemo{
		BufferMap: make(map[string]*bufio.Reader),
	}
}
