package todostore

import (
	"context"
	"testing"
	"todo/logger"
	"todo/todostore"

	"github.com/google/uuid"
)

/*
go test todo_test.go todo.go
*/
func TestTodoList(t *testing.T) {

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	t.Run("Add new todo", func(t *testing.T) {
		err := todostore.AddItemToList(ctx, "testtodolist", "Shopping", "Bread")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		list, err := todostore.GetList(ctx, "testtodolist")
		if err != nil {
			t.Errorf("Error getting todo: %s", err.Error())
		}

		if len(list.LItems) != 1 {
			t.Errorf("Expected 1 todo item, got %d", len(list.LItems))
		}

		if list.LItems[0].Name != "Shopping" {
			t.Errorf("Expected item name, got %s", list.LItems[0].Name)
		}

		_ = todostore.DeleteItemFromList(ctx, "testtodolist", "Shopping")
	})

	t.Run("Delete todo", func(t *testing.T) {
		err := todostore.AddItemToList(ctx, "testtodolist", "Shopping", "Bread")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		err = todostore.DeleteItemFromList(ctx, "testtodolist", "Shopping")
		if err != nil {
			t.Errorf("Error removing todo: %s", err.Error())
		}

		list, err := todostore.GetList(ctx, "testtodolist")
		if err != nil {
			t.Errorf("Error getting todo: %s", err.Error())
		}

		if len(list.LItems) != 0 {
			t.Errorf("Expected 0 todo item, got %d", len(list.LItems))
		}
	})

	t.Run("Update todo description", func(t *testing.T) {
		err := todostore.AddItemToList(ctx, "testtodolist", "Shopping", "Bread")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		err = todostore.UpdateListItemDescription(ctx, "testtodolist", "Shopping", "Todo Description Updated")
		if err != nil {
			t.Errorf("Error updating todo description: %s", err.Error())
		}

		list, err := todostore.GetList(ctx, "testtodolist")
		if err != nil {
			t.Errorf("Error getting todo: %s", err.Error())
		}

		if list.LItems[0].Description != "Todo Description Updated" {
			t.Errorf("Expected todo item description to be 'Todo Description Updated', got %s", list.LItems[0].Description)
		}

		_ = todostore.DeleteItemFromList(ctx, "testtodolist", "Shopping")
	})

	t.Run("Update todo status", func(t *testing.T) {
		err := todostore.AddItemToList(ctx, "testtodolist", "Shopping", "Bread")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		list, err := todostore.GetList(ctx, "testtodolist")
		if err != nil {
			t.Errorf("Error getting todo: %s", err.Error())
		}

		if list.LItems[0].Status != todostore.StatusNotStarted {
			t.Errorf("Expected todostore.StatusNotStarted, got %s", list.LItems[0].Status)
		}

		err = todostore.UpdateListItemStatus(ctx, "testtodolist", "Shopping", todostore.StatusStarted)
		if err != nil {
			t.Errorf("Error updating todo status: %s", err.Error())
		}

		if list.LItems[0].Status != todostore.StatusStarted {
			t.Errorf("Expected todostore.StatusStarted, got %s", list.LItems[0].Status)
		}

		_ = todostore.DeleteItemFromList(ctx, "testtodolist", "Shopping")
	})
}
