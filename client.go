package lftp

import (
	"fmt"
	"net"

	"github.com/Lunarisnia/lftp/internal/dsu"
	"github.com/Lunarisnia/lftp/internal/filesystem"
	"github.com/google/uuid"
)

type Client struct{}

func (c *Client) connect(address string) (dsu.LFTPConnection, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (c *Client) SendFile(address string, filePath string) error {
	connection, err := c.connect(address)
	defer connection.Close()
	if err != nil {
		return err
	}

	// TODO: A way to memoize transferred file if its bigger than x amount
	// FIXME: There should be a way to send using existing uuid later
	filesystem.OpenFile(filePath, 3)

	generatedId := uuid.NewString()
	header := dsu.LFTPHeader{
		Version:       "1.0",
		ContentLength: 1,
		TotalLength:   2,
		StartOffset:   0,
		EndOffset:     2,
		ContentID:     generatedId,
	}
	fmt.Fprintf(connection, header.ConstructString())

	return nil
}
