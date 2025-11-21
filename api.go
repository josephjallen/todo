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

func writeJSON(ctx context.Context, status int, w http.ResponseWriter, v interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, err := json.Marshal(v)
	if status >= 200 && status < 300 {
		logger.InfoLog(ctx, "Response: "+string(b))
	} else if status >= 400 {
		logger.ErrorLog(ctx, "Response: "+string(b)+" Error: "+err.Error())
	}
	logger.InfoLog(ctx, string(b))
	_ = json.NewEncoder(w).Encode(v)
}

func createListHandler(w http.ResponseWriter, r *http.Request) {

	var cr CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	_, err := todostore.CreateList(r.Context(), cr.TodoListName)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err creating list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, MessageResponse{Message: "Created List: " + cr.TodoListName}, nil)
}

func getListHandler(w http.ResponseWriter, r *http.Request) {

	var gr GetListRequest
	if err := json.NewDecoder(r.Body).Decode(&gr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	list, err := todostore.GetList(r.Context(), gr.TodoListName)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err retrieving list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusOK, w, ListResponse{List: *list}, nil)
}

func addItemHandler(w http.ResponseWriter, r *http.Request) {

	var ar AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	err := todostore.AddItemToList(r.Context(), ar.TodoListName, ar.ItemName, ar.ItemDescription)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err updating list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, MessageResponse{Message: "Added Item to List: " + ar.TodoListName}, nil)
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {

	var dr DeleteItemRequest
	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	err := todostore.DeleteItemFromList(r.Context(), dr.TodoListName, dr.ItemName)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err deleting list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, MessageResponse{Message: "Deleted Item From TodoList: " + dr.TodoListName}, nil)
}

func updateItemDescriptionHandler(w http.ResponseWriter, r *http.Request) {

	var ur UpdateItemDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	err := todostore.UpdateListItemDescription(r.Context(), ur.TodoListName, ur.ItemName, ur.ItemDescription)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err updating list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, MessageResponse{Message: "Updated TodoList: " + ur.TodoListName}, nil)
}

func updateItemStatusHandler(w http.ResponseWriter, r *http.Request) {

	var ur UpdateItemStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, MessageResponse{Message: "Invalid JSON"}, err)
		return
	}

	err := todostore.UpdateListItemStatus(r.Context(), ur.TodoListName, ur.ItemName, ur.ItemStatus)
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err updating list"}, err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, MessageResponse{Message: "Updated TodoList: " + ur.TodoListName}, nil)
}

func dynamicListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	list, err := todostore.GetList(r.Context(), "TodoList1")
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err retrieving list"}, err)
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

<h2>` + list.Name + `</h2>
<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
    <th>Status</th>
  </tr>`)
	if err != nil {
		writeJSON(r.Context(), http.StatusInternalServerError, w, MessageResponse{Message: "Internal Err"}, err)
		return
	}
	err = theader.Execute(w, nil)
	if err != nil {
		writeJSON(r.Context(), http.StatusInternalServerError, w, MessageResponse{Message: "Internal Err"}, err)
		return
	}
	/* TEMPLATING PASS IN STRUCT AND THAT SHOULD POP HTML LIST PLACEHOLDER */

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
		writeJSON(r.Context(), http.StatusInternalServerError, w, MessageResponse{Message: "Internal Err"}, err)
		return
	}
	err = tfooter.Execute(w, nil)
	if err != nil {
		writeJSON(r.Context(), http.StatusInternalServerError, w, MessageResponse{Message: "Internal Err"}, err)
		return
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
