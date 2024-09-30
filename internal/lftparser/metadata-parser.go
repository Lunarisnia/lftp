package lftparser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Lunarisnia/lftp/internal/dsu"
)

func ParseHeader(rawContent []byte) (*dsu.LFTPHeader, error) {
	content := string(rawContent)
	headers := strings.Split(content, "||||")
	if len(headers) < 8 {
		return nil, errors.New("failed to parse headers")
	}
	contentLength, err := strconv.Atoi(headers[2])
	if err != nil {
		return nil, err
	}
	totalLength, err := strconv.Atoi(headers[3])
	if err != nil {
		return nil, err
	}
	startOffset, err := strconv.Atoi(headers[4])
	if err != nil {
		return nil, err
	}
	endOffset, err := strconv.Atoi(headers[5])
	if err != nil {
		return nil, err
	}
	header := &dsu.LFTPHeader{
		Version:       headers[1],
		ContentLength: contentLength,
		TotalLength:   totalLength,
		StartOffset:   startOffset,
		EndOffset:     endOffset,
		ContentID:     headers[6],
		Content:       []byte(headers[7]),
	}
	return header, nil
}
