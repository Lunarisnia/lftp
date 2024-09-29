package filesystem

import "os"

func OpenFile(path string) error {
	// TODO: Open file and save the reference in a map
	// TODO: Return custom file struct that can consume a bytes at a time and keep track of the cursor
	_, err := os.Open(path)
	if err != nil {
		return err
	}

	return nil
}
