package filesystem

import (
	"bufio"
	"os"
)

// TODO: need custom file struct that can consume a bytes at a time and keep track of the cursor
func OpenFile(path string, bufferSize int) (*bufio.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReaderSize(f, bufferSize)

	return buf, nil
}
