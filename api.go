package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo/logger"
	"todo/todostore"

	"github.com/google/uuid"
)

var ctx context.Context

type CreateListRequest struct {
	ID           string `json:"ID"`
	TodoListName string `json:"TodoListName"`
}

type GetListRequest struct {
	ID           string `json:"ID"`
	TodoListName string `json:"TodoListName"`
}

type AddItemRequest struct {
	ID              string `json:"ID"`
	TodoListName    string `json:"TodoListName"`
	ItemName        string `json:"ItemName"`
	ItemDescription string `json:"ItemDescription"`
}

type DeleteItemRequest struct {
	ID           string `json:"ID"`
	ItemName     string `json:"ItemName"`
	TodoListName string `json:"TodoListName"`
}

type UpdateItemDescriptionRequest struct {
	ID              string `json:"ID"`
	TodoListName    string `json:"TodoListName"`
	ItemName        string `json:"ItemName"`
	ItemDescription string `json:"ItemDescription"`
}

type UpdateItemStatusRequest struct {
	ID           string `json:"ID"`
	TodoListName string `json:"TodoListName"`
	ItemName     string `json:"ItemName"`
	ItemStatus   string `json:"ItemStatus"`
}

func addTraceIDLayer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := r.Header.Get("X-Trace-ID")
		if traceID != "" {
			if ctx != nil && traceID != ctx.Value(logger.TraceIdKey{}).(string) {
				w.Header().Set("X-Trace-ID", traceID)
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Currently unable to handle multiple trace ids"})
				return
			} else if ctx == nil {
				ctx = context.WithValue(context.Background(), logger.TraceIdKey{}, traceID)
			}
		} else {
			if ctx == nil {
				ctx = context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
			}
			traceID = ctx.Value(logger.TraceIdKey{}).(string)
		}

		w.Header().Set("X-Trace-ID", traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func createListHandler(w http.ResponseWriter, r *http.Request) {

	var cr CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if cr.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, cr.TodoListName)
	if err == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "List already exists: " + cr.TodoListName})
		return
	}

	err = todostore.Init(ctx, cr.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"ok": true, "createrequest": cr})
}

func getListHandler(w http.ResponseWriter, r *http.Request) {

	var gr GetListRequest
	if err := json.NewDecoder(r.Body).Decode(&gr); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if gr.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, gr.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(ctx, gr.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"todoList": todostore.List})

}

func addItemHandler(w http.ResponseWriter, r *http.Request) {

	var ar AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if ar.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, ar.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(ctx, ar.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.AddItemToList(ctx, ar.ItemName, ar.ItemDescription)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"Add Item To TodoList": todostore.List})
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {

	var dr DeleteItemRequest
	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if dr.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, dr.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(ctx, dr.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.DeleteItemFromList(ctx, dr.ItemName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]interface{}{"Deleted Item To TodoList": todostore.List})
}

func updateItemDescriptionHandler(w http.ResponseWriter, r *http.Request) {

	var ur UpdateItemDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if ur.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, ur.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(ctx, ur.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.UpdateListItemDescription(ctx, ur.ItemName, ur.ItemDescription)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"Updated TodoList": todostore.List})
}

func updateItemStatusHandler(w http.ResponseWriter, r *http.Request) {

	var ur UpdateItemStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if ur.TodoListName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "TodoListName required"})
		return
	}

	err := todostore.CheckListExists(ctx, ur.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(ctx, ur.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.UpdateListItemStatus(ctx, ur.ItemName, ur.ItemStatus)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"Updated TodoList": todostore.List})
}

/*
go run api.go
*/
func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/createlist", createListHandler)
	mux.HandleFunc("/getlist", getListHandler)
	mux.HandleFunc("/additem", addItemHandler)
	mux.HandleFunc("/deleteitem", deleteItemHandler)
	mux.HandleFunc("/updateitemdescription", updateItemDescriptionHandler)
	mux.HandleFunc("/updateitemstatus", updateItemStatusHandler)

	handler := addTraceIDLayer(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
