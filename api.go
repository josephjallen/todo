package main

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
	"todo/logger"
	"todo/todostore"

	"github.com/google/uuid"
)

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
		if traceID == "" {
			traceID = uuid.NewString()
			w.Header().Set("X-Trace-ID", traceID)
		}

		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addLogLayer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLog(r.Context(), "Request Recieved: "+r.Method+" "+r.URL.Path)

		next.ServeHTTP(w, r)
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

	err := todostore.CheckListExists(r.Context(), cr.TodoListName)
	if err == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "List already exists: " + cr.TodoListName})
		return
	}

	err = todostore.Init(r.Context(), cr.TodoListName)

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

	err := todostore.CheckListExists(r.Context(), gr.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(r.Context(), gr.TodoListName)

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

	err := todostore.CheckListExists(r.Context(), ar.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(r.Context(), ar.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.AddItemToList(r.Context(), ar.ItemName, ar.ItemDescription)

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

	err := todostore.CheckListExists(r.Context(), dr.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(r.Context(), dr.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.DeleteItemFromList(r.Context(), dr.ItemName)

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

	err := todostore.CheckListExists(r.Context(), ur.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(r.Context(), ur.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.UpdateListItemDescription(r.Context(), ur.ItemName, ur.ItemDescription)

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

	err := todostore.CheckListExists(r.Context(), ur.TodoListName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.Init(r.Context(), ur.TodoListName)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = todostore.UpdateListItemStatus(r.Context(), ur.ItemName, ur.ItemStatus)

	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"Updated TodoList": todostore.List})
}

func dynamicListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err := todostore.Init(r.Context(), "TodoList1")
	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	var theader *template.Template
	theader, err = template.New("todolistheader").Parse(`<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2>` + todostore.List.Name + `</h2>
<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
    <th>Status</th>
  </tr>`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = theader.Execute(w, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	/*var tlistitems *template.Template
	tlistitems, err = template.New("todolistfooter").Parse(`</table></body></html>`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = tfooter.Execute(w, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}*/

	var tfooter *template.Template
	tfooter, err = template.New("todolistfooter").Parse(`</table></body></html>`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = tfooter.Execute(w, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

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
	mux.HandleFunc("/list", dynamicListHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.ToLower(r.URL.Path), ".html") {
			http.ServeFile(w, r, "./static/"+path.Base(r.URL.Path))
		} else {
			http.ServeFile(w, r, "./static/"+path.Base(r.URL.Path)+".html")
		}
	})

	handler := addLogLayer(mux)
	handler = addTraceIDLayer(handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Parallel goroutine to run the server
	go func() {
		logger.InfoLog(nil, "Server listening on :8080")
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorLog(nil, "HTTP server Close error: "+err.Error())
		}
		logger.InfoLog(nil, "Server stopped listening on :8080")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.ErrorLog(nil, "HTTP shutdown error: "+err.Error())
		err = srv.Close()
		if err != nil {
			logger.ErrorLog(nil, "Server shutdown error: "+err.Error())
		}
	}
	logger.InfoLog(nil, "Shutdown complete")
}
