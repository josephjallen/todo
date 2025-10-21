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
go run todo.go -todoList=todod1 -deleteitemname=monday
go run todo.go -todoList=todod1 -updateitemname=monday -updateitemdescription=gotoshop_updated
*/
func main() {

	todoListName := flag.String("todoList", "", "")
	todoAddItemName := flag.String("additemname", "", "")
	todoAddItemDescription := flag.String("additemdescription", "", "")

	todoUpdateItemName := flag.String("updateitemname", "", "")
	todoUpdateItemDescription := flag.String("updateitemdescription", "", "")

	todoDeleteItemName := flag.String("deleteitemname", "", "")

	flag.Parse()

	list, err := getList(*todoListName)
	if err != nil {
		fmt.Println("error:", err)
	}

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		for _, lItem := range list.LItems {
			if lItem.Name == *todoAddItemName {
				fmt.Println("Item already exists: " + lItem.Name)
				return
			}
		}
		lItem := todoListItem{Name: *todoAddItemName, Description: *todoAddItemDescription}
		list.LItems = append(list.LItems, lItem)
		fmt.Println("Added item: " + lItem.Name + " to list: " + list.Name)
	} else if *todoUpdateItemName != "" && *todoUpdateItemDescription != "" {
		var updateItemIndex int = -2
		for index, lItem := range list.LItems {
			if lItem.Name == *todoUpdateItemName {
				updateItemIndex = index
				break
			}
		}
		if updateItemIndex > -2 {
			list.LItems[updateItemIndex].Description = *todoUpdateItemDescription
			fmt.Println("Item Updated: " + *todoUpdateItemName + " in list: " + list.Name)
		} else {
			fmt.Println("Cannot find Item to update: " + *todoUpdateItemName)
		}
	} else if *todoDeleteItemName != "" {
		var deleteItemIndex int = -2
		for index, lItem := range list.LItems {
			if lItem.Name == *todoDeleteItemName {
				deleteItemIndex = index
				break
			}
		}
		if deleteItemIndex > -2 {
			list.LItems = append(list.LItems[:deleteItemIndex], list.LItems[deleteItemIndex+1:]...)
			fmt.Println("Item Deleted: " + *todoDeleteItemName + " from list: " + list.Name)
		} else {
			fmt.Println("Cannot find Item to delete: " + *todoDeleteItemName)
		}
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
