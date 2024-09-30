package lftp

import (
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
			writeCounter := 0
			buf, err := os.Create("./tests/32bytes-copy")
			assert.Nil(t, err)

			server = NewLFTPServer(func(header *dsu.LFTPHeader) {
				n, err := buf.WriteString(string(header.Content))
				assert.Nil(t, err)
				writeCounter += n
				if writeCounter == 32 {
					waiting <- true
					buf.Close()
				}
			})
			server.Listen(":6968")
		}()
		client := NewLFTPClient()
		err := client.SendFile("localhost:6968", "./tests/32bytes", 16)
		assert.Nil(t, err)
		<-waiting
		err = server.Close()
		assert.Nil(t, err)
	})
}
