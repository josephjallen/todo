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

	err := todostore.Init(*todoListName, ctx)
	if err != nil {
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
	}

	if *todoAddItemName != "" && *todoAddItemDescription != "" {
		todostore.AddItemToList(*todoAddItemName, *todoAddItemDescription, ctx)
	} else if *todoUpdateItemName != "" && *todoUpdateItemDescription != "" {
		todostore.UpdateListItem(*todoUpdateItemName, *todoUpdateItemDescription, ctx)
	} else if *todoDeleteItemName != "" {
		todostore.DeleteItemFromList(*todoDeleteItemName, ctx)
	}

	err = todostore.SaveList(ctx)
	if err != nil {
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
	}

}
