package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
	"todo/logger"
)

type FileStorage struct {
	Name string
}

var fileStorage *FileStorage

func Init(ctx context.Context, fileName string) error {
	if fileStorage == nil {
		logger.InfoLog(ctx, "FileStorage Creating single instance now.")
		fileStorage = &FileStorage{Name: fileName}
	} else {
		logger.WarningLog(ctx, "FileStorage Single instance already created.")
	}

	return nil
}

func SaveByteSliceToFile(val []byte) error {

	err := backupFile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileStorage.Name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return (err)
	}
	defer f.Close()

	f.Write(val)

	return nil
}

func LoadFileToByteSlice(ctx context.Context) ([]byte, error) {
	_, err := os.Stat(fileStorage.Name)
	if os.IsNotExist(err) {
		logger.InfoLog(ctx, "File does not exist: "+fileStorage.Name)
		return nil, nil
	}

	f, err := os.OpenFile(fileStorage.Name, os.O_RDONLY, 0644)
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
			b = append(b, buf[:n]...)
		}
	}

	return b, nil
}

func backupFile() error {
	if _, err := os.Stat(fileStorage.Name); err == nil {
		// Open the source file
		sourceFile, err := os.Open(fileStorage.Name)
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer sourceFile.Close()

		// Create the destination file
		destinationFile, err := os.Create("./Backups/" + fileStorage.Name + "_" + time.Now().String())
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destinationFile.Close()

		// Copy the content
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		err = os.Remove(fileStorage.Name)
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}

	}

	return nil
}
