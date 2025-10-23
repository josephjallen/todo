package main

import (
	"flag"
	"todo/logger"
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
		logger.ErrorLog.Println("error:", err)
	}

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		todostore.AddItemToList(list, *todoAddItemName, *todoAddItemDescription)
	} else if *todoUpdateItemName != "" && *todoUpdateItemDescription != "" {
		todostore.UpdateListItem(list, *todoUpdateItemName, *todoUpdateItemDescription)
	} else if *todoDeleteItemName != "" {
		todostore.DeleteItemFromList(list, *todoDeleteItemName)
	}

	err = todostore.SaveList(list)
	if err != nil {
		logger.ErrorLog.Println("error:", err)
	}

}
