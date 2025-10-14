package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"todo/filestorage"
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

	b, _ := filestorage.LoadFileToByteSlice("file")
	fmt.Println(string(b))

	todoListName := flag.String("todoList", "foo", "a string")
	todoItemName := flag.String("item", "42", "a string")

	flag.Parse()

	list := todoList{Name: *todoListName}
	lItem := todoListItem{Name: *todoItemName}

	list.addItem(lItem)

	bb, _ := json.Marshal(list)

	fmt.Println(string(bb))

	filestorage.SaveByteSliceToFile(bb)
}
