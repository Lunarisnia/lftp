package filesystem

import (
	"bufio"
	"io"
	"os"
)

// TODO: Need to somehow save this opened file reader in a map
// TODO: need custom file struct that can consume a bytes at a time and keep track of the cursor
func OpenFile(path string, bufferSize int) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReaderSize(f, bufferSize)

	return buf, nil
}
