package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"todo/filestorage"
)

type todoList struct {
	Name   string         `json:"name"`
	LItems []todoListItem `json:"litems"`
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

	todoListName := flag.String("todoList", "todo", "a string")
	todoItemName := flag.String("item", "1", "a string")

	flag.Parse()

	list, err := getList(*todoListName)
	if err != nil {
		fmt.Println("error:", err)
	}

	if todoItemName != nil {
		addItem(list, *todoItemName)
	}

	list_bb, err := json.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}

	filestorage.SaveByteSliceToFile(list_bb)

	fmt.Println(string(list_bb))

}

func addItem(list *todoList, todoItemName string) {
	lItem := todoListItem{Name: todoItemName}

	list.addItem(lItem)
}

/* go run todo.go -list=opt -item=7
 */
func getList(todoListName string) (*todoList, error) {

	filename := todoListName
	filename += ".json"
	fmt.Println(string(filename))

	list_b, err := filestorage.LoadFileToByteSlice(filename)
	if err != nil {
		fmt.Println("error:", err)
	}

	var list todoList

	if list_b != nil {
		fmt.Println(string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		list = todoList{Name: todoListName}
	}

	return &list, nil
}
