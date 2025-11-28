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
go run cli.go -todoList=todod1 -updateitemname=monday -updateitemstatus=started
*/
func main() {

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	todoListName := flag.String("todoList", "", "")
	todoAddItemName := flag.String("additemname", "", "")
	todoAddItemDescription := flag.String("additemdescription", "", "")

	todoUpdateItemName := flag.String("updateitemname", "", "")
	todoUpdateItemDescription := flag.String("updateitemdescription", "", "")
	todoUpdateItemStatus := flag.String("updateitemstatus", "", "")

	todoDeleteItemName := flag.String("deleteitemname", "", "")

	flag.Parse()

	var err error

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		err = todostore.AddItemToList(ctx, *todoListName, *todoAddItemName, *todoAddItemDescription)
	} else if *todoUpdateItemName != "" && *todoUpdateItemDescription != "" {
		err = todostore.UpdateListItemDescription(ctx, *todoListName, *todoUpdateItemName, *todoUpdateItemDescription)
	} else if *todoUpdateItemName != "" && *todoUpdateItemStatus != "" {
		err = todostore.UpdateListItemStatus(ctx, *todoListName, *todoUpdateItemName, *todoUpdateItemStatus)
	} else if *todoDeleteItemName != "" {
		err = todostore.DeleteItemFromList(ctx, *todoListName, *todoDeleteItemName)
	}

	if err == nil {
		todostore.SaveList(ctx, *todoListName)
	}

	if err != nil {
		logger.GetCtxLogger(ctx).Error(err.Error())
	}
}
