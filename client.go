package lftp

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/Lunarisnia/lftp/internal/dsu"
	"github.com/Lunarisnia/lftp/internal/filesystem"
	"github.com/google/uuid"
)

const (
	ClientVersion = "1.0"
)

type LFTPClient interface {
	SendFile(address string, filePath string, chunkSize int) error
}

type Client struct {
	memo *dsu.ClientMemo
}

func NewLFTPClient() LFTPClient {
	return &Client{
		memo: dsu.NewClientMemo(),
	}
}

func (c *Client) connect(address string) (dsu.LFTPConnection, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (c *Client) SendFile(address string, filePath string, chunkSize int) error {
	generatedId := uuid.NewString()
	if _, exist := c.memo.BufferMap[generatedId]; exist {
		fmt.Printf("The impossible just happened. an UUID just clashed: %v\n", generatedId)
	}
	f, err := filesystem.OpenFile(filePath, chunkSize)
	if err != nil {
		return err
	}

	for {
		connection, err := c.connect(address)
		defer connection.Close()
		if err != nil {
			return err
		}

		c.memo.BufferMap[generatedId] = f
		c.memo.TotalLength = f.Size() // FIXME: This does not return the entire filesize this just return the ucrrent buffer size
		header := dsu.LFTPHeader{
			Version:       ClientVersion,
			ContentLength: f.Size(),
			TotalLength:   c.memo.TotalLength,
			StartOffset:   0,
			EndOffset:     chunkSize,
			ContentID:     generatedId,
			Content:       make([]byte, chunkSize),
		}
		n, err := f.Read(header.Content)
		fmt.Printf("%v amount of data has been read.\n", n)
		if err != nil {
			return err
		}

		response := header.ConstructString()
		fmt.Fprint(connection, response)

		preview, err := f.Peek(chunkSize)
		if err != nil {
			if errors.Is(err, io.EOF) && len(preview) == 0 {
				return nil
			}
			if !errors.Is(err, io.EOF) {
				return err
			}
		}
	}
}
