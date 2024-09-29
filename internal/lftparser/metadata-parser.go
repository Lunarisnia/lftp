package lftparser

import (
	"fmt"

	"github.com/Lunarisnia/lftp/internal/dsu"
)

func ParseHeader(rawContent []byte) *dsu.LFTPHeader {
	content := string(rawContent)
	fmt.Println(content)
	header := &dsu.LFTPHeader{}
	return header
}
