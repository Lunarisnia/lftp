package lftp

import (
	"bytes"
	"os"
	"testing"

	"github.com/Lunarisnia/lftp/internal/dsu"
	"github.com/stretchr/testify/assert"
)

func Test_SendingFile(t *testing.T) {
	t.Run("Expect to successfully send the entire file", func(t *testing.T) {
		var server LFTPServer
		waiting := make(chan bool)
		go func() {
			server = NewLFTPServer(func(header *dsu.LFTPHeader) {
				waiting <- true
			})
			server.Listen(":6969")
		}()
		client := NewLFTPClient()
		err := client.SendFile("localhost:6969", "./tests/7bytes", 16)
		assert.Nil(t, err)
		<-waiting
		err = server.Close()
		assert.Nil(t, err)
	})
	t.Run("Expect to successfully send the entire file step by step", func(t *testing.T) {
		var server LFTPServer
		waiting := make(chan bool)
		go func() {
			buf, err := os.Create("./tests/ordinary-copy")
			assert.Nil(t, err)

			server = NewLFTPServer(func(header *dsu.LFTPHeader) {
				nonZero := header.Content
				if header.ContentLength == -1 {
					nonZero = bytes.Trim(
						header.Content,
						"\x00",
					) // NOTE: This might cause error in files that intentionally pad with zero values
				}
				_, err := buf.Write(nonZero)
				assert.Nil(t, err)
				if header.ContentLength == -1 {
					waiting <- true
					buf.Close()
				}
			})
			server.Listen(":6968")
		}()
		client := NewLFTPClient()
		err := client.SendFile("localhost:6968", "./tests/ordinary", 1000)
		assert.Nil(t, err)
		<-waiting
		err = server.Close()
		assert.Nil(t, err)
	})
}
