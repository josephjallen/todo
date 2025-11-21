package web

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"todo/logger"
	"todo/todostore"

	"github.com/google/uuid"
)

func AddTraceIDLayer(next http.Handler) http.Handler {
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

func AddLogLayer(next http.Handler) http.Handler {
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
	_ = json.NewEncoder(w).Encode(v)
}

func CreateListHandler(w http.ResponseWriter, r *http.Request) {

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

func GetListHandler(w http.ResponseWriter, r *http.Request) {

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

func AddItemHandler(w http.ResponseWriter, r *http.Request) {

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

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {

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

func UpdateItemDescriptionHandler(w http.ResponseWriter, r *http.Request) {

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

func UpdateItemStatusHandler(w http.ResponseWriter, r *http.Request) {

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

func DynamicListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	list, err := todostore.GetList(r.Context(), r.PathValue("listname"))
	if err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, MessageResponse{Message: "Err retrieving list"}, err)
		return
	}

	tmpl := template.Must(template.ParseFiles("./web/dynamic/list.html"))
	tmpl.Execute(w, list)
}
