package todostore

import (
	"encoding/json"
	"fmt"
	"todo/filestorage"
)

type TodoList struct {
	Name   string         `json:"name"`
	LItems []TodoListItem `json:"litems"`
}

type TodoListItem struct {
	Name        string
	Description string
}

func AddItemToList(list *TodoList, itemName string, itemDescription string) {
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			fmt.Println("Item already exists: " + lItem.Name)
			return
		}
	}
	lItem := TodoListItem{Name: itemName, Description: itemDescription}
	list.LItems = append(list.LItems, lItem)
	fmt.Println("Added item: " + lItem.Name + " to list: " + list.Name)
}

/*
 */
func GetList(todoListName string) (*TodoList, error) {

	filename := todoListName
	filename += ".json"
	fmt.Println(string(filename))

	list_b, err := filestorage.LoadFileToByteSlice(filename)
	if err != nil {
		fmt.Println("error:", err)
	}

	var list TodoList

	if list_b != nil {
		fmt.Println(string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		list = TodoList{Name: todoListName}
	}

	return &list, nil
}
