package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"todo/filestorage"
	"todo/todostore"
)

/*
go run cli.go -todoList=todod1 -additemname=monday -additemdescription=gotoshop
go run cli.go -todoList=todod1 -deleteitemname=monday
go run cli.go -todoList=todod1 -updateitemname=monday -updateitemdescription=gotoshop_updated
*/
func main() {

	todoListName := flag.String("todoList", "", "")
	todoAddItemName := flag.String("additemname", "", "")
	todoAddItemDescription := flag.String("additemdescription", "", "")

	todoUpdateItemName := flag.String("updateitemname", "", "")
	todoUpdateItemDescription := flag.String("updateitemdescription", "", "")

	todoDeleteItemName := flag.String("deleteitemname", "", "")

	flag.Parse()

	list, err := todostore.GetList(*todoListName)
	if err != nil {
		fmt.Println("error:", err)
	}

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		todostore.AddItemToList(list, *todoAddItemName, *todoAddItemDescription)
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
