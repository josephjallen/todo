package main

import (
	"flag"
	"todo/logger"
	"todo/todostore"

	"context"

	"github.com/google/uuid"
)

/*
go run cli.go -todoList=todod1 -additemname=monday -additemdescription=gotoshop
go run cli.go -todoList=todod1 -deleteitemname=monday
go run cli.go -todoList=todod1 -updateitemname=monday -updateitemdescription=gotoshop_updated
*/
func main() {

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	todoListName := flag.String("todoList", "", "")
	todoAddItemName := flag.String("additemname", "", "")
	todoAddItemDescription := flag.String("additemdescription", "", "")

	todoUpdateItemName := flag.String("updateitemname", "", "")
	todoUpdateItemDescription := flag.String("updateitemdescription", "", "")

	todoDeleteItemName := flag.String("deleteitemname", "", "")

	flag.Parse()

	err := todostore.Init(ctx, *todoListName)

	if err != nil {
		logger.ErrorLog(ctx, err.Error())
	}

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		err = todostore.AddItemToList(ctx, *todoAddItemName, *todoAddItemDescription)
	} else if *todoUpdateItemName != "" && *todoUpdateItemDescription != "" {
		err = todostore.UpdateListItem(ctx, *todoUpdateItemName, *todoUpdateItemDescription)
	} else if *todoDeleteItemName != "" {
		err = todostore.DeleteItemFromList(ctx, *todoDeleteItemName)
	}

	if err != nil {
		logger.ErrorLog(ctx, err.Error())
	}
}
