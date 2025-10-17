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
	Name        string
	Description string
}

/*
go run todo.go -todoList=todod1 -additemname=monday -additemdescription=gotoshop
*/
func main() {

	todoListName := flag.String("todoList", "todo", "a list name")
	todoAddItemName := flag.String("additemname", "1", "a list item name")
	todoAddItemDescription := flag.String("additemdescription", "1", "a description")

	flag.Parse()

	list, err := getList(*todoListName)
	if err != nil {
		fmt.Println("error:", err)
	}

	if todoAddItemName != nil && todoAddItemDescription != nil {
		for _, lItem := range list.LItems {
			if lItem.Name == *todoAddItemName {
				fmt.Println("Item already exists: " + lItem.Name)
				return
			}
		}
		lItem := todoListItem{Name: *todoAddItemName, Description: *todoAddItemDescription}
		list.LItems = append(list.LItems, lItem)
		fmt.Println("Added item: " + lItem.Name + " to list: " + list.Name)
	}

	list_bb, err := json.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}

	filestorage.SaveByteSliceToFile(list_bb, *todoListName+".json")

	fmt.Println(string(list_bb))

}

/*
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
