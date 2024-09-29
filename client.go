package lftp

import (
	"errors"
	"fmt"
	"net"

	"github.com/Lunarisnia/lftp/internal/dsu"
)

type Client struct {
	connection dsu.LFTPConnection
}

func (c *Client) connect(address string) error {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.connection = connection

	return nil
}

func (c *Client) SendRequest(address string) error {
	c.connection = nil

	c.connect(address)
	if c.connection == nil {
		return errors.New("failed to connect")
	}

	header := dsu.LFTPHeader{
		Version:       "1.0",
		ContentLength: 1,
		TotalLength:   2,
		StartOffset:   0,
		EndOffset:     2,
		ContentID:     "UUIDHERE",
	}
	fmt.Fprintf(c.connection, header.ConstructString())

	defer func() {
		c.connection.Close()
		c.connection = nil
	}()
	return nil
}
