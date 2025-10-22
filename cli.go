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
		todostore.UpdateListItem(list, *todoUpdateItemName, *todoUpdateItemDescription)
	} else if *todoDeleteItemName != "" {
		todostore.DeleteItemFromList(list, *todoDeleteItemName)
	}

	list_bb, err := json.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}

	filestorage.SaveByteSliceToFile(list_bb, *todoListName+".json")

	fmt.Println(string(list_bb))

}
