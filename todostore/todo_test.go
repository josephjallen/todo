package todostore

import (
	"context"
	"testing"
	"todo/logger"
	"todo/todostore"

	"github.com/google/uuid"
)

func initTodoStore() {
	todostore.List = &todostore.TodoList{Name: "testtodolist"}
}

/*
go test todo_test.go todo.go
*/
func TestTodoList(t *testing.T) {

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	t.Run("Check name", func(t *testing.T) {
		initTodoStore()
		if todostore.List.Name != "testtodolist" {
			t.Errorf("Expected list name to be 'testtodolist', got '%s'", todostore.List.Name)
		}
	})

	t.Run("Add new todo", func(t *testing.T) {
		initTodoStore()
		err := todostore.AddItemToList(ctx, "testtodolist", "Todo 1 Description")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		if len(todostore.List.LItems) != 1 {
			t.Errorf("Expected 1 todo item, got %d", len(todostore.List.LItems))
		}
	})

	t.Run("Delete todo", func(t *testing.T) {
		initTodoStore()
		err := todostore.AddItemToList(ctx, "testtodolist", "Todo 1 Description")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		if len(todostore.List.LItems) != 1 {
			t.Errorf("Expected 1 todo item, got %d", len(todostore.List.LItems))
		}

		err = todostore.DeleteItemFromList(ctx, "testtodolist")
		if err != nil {
			t.Errorf("Error removing todo: %s", err.Error())
		}

		if len(todostore.List.LItems) != 0 {
			t.Errorf("Expected 0 todo item, got %d", len(todostore.List.LItems))
		}
	})

	t.Run("Update todo description", func(t *testing.T) {
		initTodoStore()
		err := todostore.AddItemToList(ctx, "testtodolist", "Todo 1 Description")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		if len(todostore.List.LItems) != 1 {
			t.Errorf("Expected 1 todo item, got %d", len(todostore.List.LItems))
		}

		err = todostore.UpdateListItemDescription(ctx, "testtodolist", "Todo 1 Description Updated")
		if err != nil {
			t.Errorf("Error updating todo description: %s", err.Error())
		}

		if todostore.List.LItems[0].Description != "Todo 1 Description Updated" {
			t.Errorf("Expected todo item description to be 'Todo 1 Description Updated', got %s", todostore.List.LItems[0].Description)
		}
	})

	t.Run("Update todo status", func(t *testing.T) {
		initTodoStore()
		err := todostore.AddItemToList(ctx, "testtodolist", "Todo 1 Description")
		if err != nil {
			t.Errorf("Error adding new todo: %s", err.Error())
		}

		if len(todostore.List.LItems) != 1 {
			t.Errorf("Expected 1 todo item, got %d", len(todostore.List.LItems))
		}

		if todostore.List.LItems[0].Status != todostore.StatusNotStarted {
			t.Errorf("Expected todostore.StatusNotStarted, got %s", todostore.List.LItems[0].Status)
		}

		err = todostore.UpdateListItemStatus(ctx, "testtodolist", todostore.StatusStarted)
		if err != nil {
			t.Errorf("Error updating todo status: %s", err.Error())
		}

		if todostore.List.LItems[0].Status != todostore.StatusStarted {
			t.Errorf("Expected todostore.StatusStarted, got %s", todostore.List.LItems[0].Status)
		}
	})
}
