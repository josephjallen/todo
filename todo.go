package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type todoList struct {
	Name   string `json:"name"`
	LItems []todoListItem
}

type todoListItem struct {
	Name string
}

func (l *todoList) addItem(lItem todoListItem) {
	l.LItems = append(l.LItems, lItem)
	fmt.Println("Added item: " + lItem.Name + " to list: " + l.Name)
}

/* go run todo.go -list=opt -item=7
 */
func main() {

	b, _ := loadFileToByteSlice("file")
	fmt.Println(string(b))

	todoListName := flag.String("todoList", "foo", "a string")
	todoItemName := flag.String("item", "42", "a string")

	flag.Parse()

	list := todoList{Name: *todoListName}
	lItem := todoListItem{Name: *todoItemName}

	list.addItem(lItem)

	bb, _ := json.Marshal(list)

	fmt.Println(string(bb))

	saveByteSliceToFile(bb)
}

func saveByteSliceToFile(val []byte) error {

	if _, err := os.Stat("File"); err == nil {
		// Open the source file
		sourceFile, err := os.Open("file")
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer sourceFile.Close()

		// Create the destination file
		destinationFile, err := os.Create("./Backups/file_" + time.Now().String())
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destinationFile.Close()

		// Copy the content
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		err = os.Remove("File")
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}

	}

	f, err := os.OpenFile("file", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return (err)
	}
	defer f.Close()

	f.Write(val)

	return nil
}

func loadFileToByteSlice(file string) ([]byte, error) {
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
