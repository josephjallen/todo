package todostore

import (
	"testing"
)

func CreateTodo() *TodoList {
	list := &TodoList{Name: "testtodolist"}
	list.LItems = []TodoListItem{
		{Name: "Item1", Description: "Todo 1"},
		{Name: "Item2", Description: "Todo 2"},
		{Name: "Item3", Description: "Todo 3"},
	}

	return list
}

/*
go test todo_test.go todo.go
*/
func TestTodoList(t *testing.T) {
	list := CreateTodo()

	t.Run("Check todolist name", func(t *testing.T) {

		if list.Name != "testtodolist" {
			t.Errorf("Expected 'testtodolist', got '%s'", list.Name)
		}
	})

	t.Run("Retrieve existing todo item", func(t *testing.T) {

		listItem := list.LItems[0]

		if listItem.Description != "Todo 1" {
			t.Errorf("Expected 'Todo 1', got '%s'", listItem.Description)
		}
	})
	/*
		t.Run("Add new todo", func(t *testing.T) {
			newTodo := types.NewTodo("New Todo", nil)
			id, _ := store.AddTodo(ctx, newTodo)

			todo, err := store.GetTodo(ctx, id)
			if err != nil {
				t.Fatalf("Expected to retrieve newly added todo with ID %s, got error: %v", id, err)
			}
			if todo.Description != "New Todo" {
				t.Errorf("Expected description 'New Todo', got '%s'", todo.Description)
			}
		})

		t.Run("Update todo status", func(t *testing.T) {
			err := store.UpdateTodoStatus(ctx, id1, types.Completed)
			if err != nil {
				t.Fatalf("Expected to update status of todo with ID %s, got error: %v", id1, err)
			}

			todo, err := store.GetTodo(ctx, id1)
			if err != nil {
				t.Fatalf("Expected to retrieve todo with ID %s, got error: %v", id1, err)
			}
			if todo.Status != types.Completed {
				t.Errorf("Expected status 'Completed', got '%s'", todo.Status)
			}
		})

		t.Run("Get todos by status", func(t *testing.T) {
			todos := store.GetTodosByStatus(ctx, types.NotStarted)
			if len(todos) != 2 {
				t.Errorf("Expected 2 not started todos, got %d", len(todos))
			}
		})

		t.Run("Get all todos", func(t *testing.T) {
			todos := store.GetAllTodos(ctx)
			if len(todos) != 2 {
				t.Errorf("Expected 2 todos (excluding completed), got %d", len(todos))
			}
		})

		t.Run("Retrieve non-existent todo", func(t *testing.T) {
			_, err := store.GetTodo(ctx, "non-existent-id")
			if err == nil {
				t.Fatalf("Expected error when retrieving non-existent todo, got nil")
			}
		})
	*/
}
