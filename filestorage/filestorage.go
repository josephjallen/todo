package filestorage

import (
	"fmt"
	"io"
	"os"
	"time"
)

func SaveByteSliceToFile(val []byte) error {

	err := backupFile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile("todo.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return (err)
	}
	defer f.Close()

	f.Write(val)

	return nil
}

func LoadFileToByteSlice(file string) ([]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := []byte{}
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			//fmt.Println(string(buf[:n]))
			b = append(b, buf...)
		}
	}

	return b, nil
}

func backupFile() error {
	if _, err := os.Stat("todo.json"); err == nil {
		// Open the source file
		sourceFile, err := os.Open("todo.json")
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer sourceFile.Close()

		// Create the destination file
		destinationFile, err := os.Create("./Backups/todo.json_" + time.Now().String())
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destinationFile.Close()

		// Copy the content
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		err = os.Remove("todo.json")
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}

	}

	return nil
}
