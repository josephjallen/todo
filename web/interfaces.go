package web

import "todo/todostore"

type CreateListRequest struct {
	TodoListName string `json:"TodoListName"`
}

type GetListRequest struct {
	TodoListName string `json:"TodoListName"`
}

type AddItemRequest struct {
	TodoListName    string `json:"TodoListName"`
	ItemName        string `json:"ItemName"`
	ItemDescription string `json:"ItemDescription"`
}

type DeleteItemRequest struct {
	ItemName     string `json:"ItemName"`
	TodoListName string `json:"TodoListName"`
}

type UpdateItemDescriptionRequest struct {
	TodoListName    string `json:"TodoListName"`
	ItemName        string `json:"ItemName"`
	ItemDescription string `json:"ItemDescription"`
}

type UpdateItemStatusRequest struct {
	TodoListName string `json:"TodoListName"`
	ItemName     string `json:"ItemName"`
	ItemStatus   string `json:"ItemStatus"`
}

type MessageResponse struct {
	Message string `json:"Message"`
}

type ListResponse struct {
	List todostore.TodoList `json:"TodoList"`
}
