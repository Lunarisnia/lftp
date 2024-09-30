package dsu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConstructingLFTPHeader(t *testing.T) {
	t.Run("All field successfully constructed", func(t *testing.T) {
		header := LFTPHeader{
			Version:       "1.0",
			ContentID:     "CONTENTIDHERE",
			StartOffset:   1,
			EndOffset:     2,
			ContentLength: 3,
			TotalLength:   10,
		}
		headerString := header.ConstructString()
		assert.Equal(t, "LFTP||||1.0||||3||||10||||1||||2||||CONTENTIDHERE||||", headerString)
	})
}
