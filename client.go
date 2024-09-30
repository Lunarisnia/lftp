package lftp

import (
	"fmt"
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
	memo dsu.ClientMemo
}

func NewClient() LFTPClient {
	return &Client{
		memo: make(dsu.ClientMemo),
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
	connection, err := c.connect(address)
	defer connection.Close()
	if err != nil {
		return err
	}

	// TODO: A way to memoize transferred file if its bigger than x amount
	// FIXME: There should be a way to send using existing uuid later
	f, err := filesystem.OpenFile(filePath, chunkSize)
	if err != nil {
		return err
	}

	generatedId := uuid.NewString()
	if _, exist := c.memo[generatedId]; exist {
		fmt.Printf("The impossible just happened. an UUID just clashed: %v\n", generatedId)
	}
	if chunkSize < f.Size() {
		c.memo[generatedId] = f
	}
	header := dsu.LFTPHeader{
		Version:       "1.0",
		ContentLength: chunkSize,
		TotalLength:   f.Size(),
		StartOffset:   0,
		EndOffset:     chunkSize,
		ContentID:     generatedId,
	}
	fmt.Fprint(connection, header.ConstructString())
	// TODO: need to call continueSendFile if its not done yet

	return nil
}

func (c *Client) continueSendFile(i int, address string, contentId string, chunkSize int) error {
	connection, err := c.connect(address)
	defer connection.Close()
	if err != nil {
		return err
	}

	f, exist := c.memo[contentId]
	if !exist {
		fmt.Println("File did not exist in the memo")
		return nil
	}

	header := dsu.LFTPHeader{
		Version:       ClientVersion,
		ContentLength: chunkSize,
		TotalLength:   f.Size(),
		StartOffset:   chunkSize,
		EndOffset:     i * chunkSize,
		ContentID:     contentId,
	}
	fmt.Fprint(connection, header.ConstructString())
	return nil
}
