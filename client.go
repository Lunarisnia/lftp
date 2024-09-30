package lftp

import (
	"bufio"
	"errors"
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
	connection, err := c.connect(address)
	defer connection.Close()
	if err != nil {
		return err
	}

	f, err := filesystem.OpenFile(filePath, chunkSize)
	if err != nil {
		return err
	}

	generatedId := uuid.NewString()
	if _, exist := c.memo.BufferMap[generatedId]; exist {
		fmt.Printf("The impossible just happened. an UUID just clashed: %v\n", generatedId)
	}
	if chunkSize < f.Size() {
		c.memo.BufferMap[generatedId] = f
	}
	c.memo.TotalLength = f.Size()
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
	// TODO: need to call continueSendFile if its not done yet
	// Peek the next chunkSize bytes

	// preview, err := f.Peek(chunkSize)
	// fmt.Println("1=======================================================")
	// if err != nil {
	// 	if !errors.Is(err, bufio.ErrBufferFull) {
	// 		return err
	// 	}
	// 	if !errors.Is(err, io.EOF) {
	// 		return err
	// 	}
	// 	delete(c.memo.BufferMap, header.ContentID)
	// 	return err
	// }
	// if err != nil && errors.Is(err, bufio.ErrBufferFull) && len(preview) == 0 {
	// 	fmt.Println("3=======================================================")
	// 	fmt.Println("THIS RAN")
	// 	delete(c.memo.BufferMap, header.ContentID)
	// 	return nil
	// }
	// c.continueSendFile(1, address, header.ContentID, chunkSize)

	return nil
}

// FIXME: Maybe just do it iteratively?
func (c *Client) continueSendFile(i int, address string, contentId string, chunkSize int) error {
	connection, err := c.connect(address)
	defer connection.Close()
	if err != nil {
		return err
	}

	f, exist := c.memo.BufferMap[contentId]
	if !exist {
		fmt.Println("File did not exist in the memo")
		return nil
	}

	header := dsu.LFTPHeader{
		Version:       ClientVersion,
		ContentLength: f.Size(),
		TotalLength:   c.memo.TotalLength,
		StartOffset:   i * chunkSize,
		EndOffset:     (i + 1) * chunkSize,
		ContentID:     contentId,
	}
	fmt.Fprint(connection, header.ConstructString())
	preview, err := f.Peek(chunkSize)
	if err != nil && !errors.Is(err, bufio.ErrBufferFull) {
		return err
	}
	if err != nil && errors.Is(err, bufio.ErrBufferFull) && len(preview) == 0 {
		fmt.Println("THIS RAN")
		delete(c.memo.BufferMap, header.ContentID)
		return nil
	}
	c.continueSendFile(i+1, address, contentId, chunkSize)
	return nil
}
