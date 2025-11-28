package web

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"todo/actors"
	"todo/logger"

	"github.com/google/uuid"
)

type Response struct {
	Message string
}

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
		logger.GetCtxLogger(r.Context()).Info("Request Recieved: " + r.Method + " " + r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func writeJSON(ctx context.Context, status int, w http.ResponseWriter, v interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(v)
	if status >= 200 && status < 300 {
		logger.GetCtxLogger(ctx).Info("Response: " + string(b))
	} else if status >= 400 {
		if err != nil {
			logger.GetCtxLogger(ctx).Error("Response: " + string(b) + " Error: " + err.Error())
		} else {
			logger.GetCtxLogger(ctx).Error("Response: " + string(b))
		}
	}
	_ = json.NewEncoder(w).Encode(v)
}

func callActor(r *http.Request, request actors.Request) chan (actors.Response) {
	responseChannel := make(chan actors.Response)

	actors.GetActor().SendMessage(r.Context(), actors.Message{
		Request:      request,
		ResponseChan: responseChannel,
		Ctx:          r.Context(),
	})

	return responseChannel
}

func CreateListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodPost}, nil)
		return
	}
	var cr CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "CreateList", TodoListName: cr.TodoListName})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err creating list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, Response{Message: "Created List: " + cr.TodoListName}, nil)
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodGet}, nil)
		return
	}
	var gr GetListRequest
	if err := json.NewDecoder(r.Body).Decode(&gr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "GetList", TodoListName: gr.TodoListName})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err retrieving list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusOK, w, response.List, nil)
}

func AddItemHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodPost}, nil)
		return
	}
	var ar AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "AddItem", TodoListName: ar.TodoListName, ItemName: ar.ItemName, ItemDescription: ar.ItemDescription})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err updating list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, Response{Message: "Added Item to List: " + ar.TodoListName}, nil)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodPost}, nil)
		return
	}
	var dr DeleteItemRequest
	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "DeleteItem", TodoListName: dr.TodoListName, ItemName: dr.ItemName})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err deleting list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, Response{Message: "Deleted Item From TodoList: " + dr.TodoListName}, nil)
}

func UpdateItemDescriptionHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodPatch}, nil)
		return
	}
	var ur UpdateItemDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "UpdateItemDescription", TodoListName: ur.TodoListName, ItemName: ur.ItemName, ItemDescription: ur.ItemDescription})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err updating list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, Response{Message: "Updated TodoList: " + ur.TodoListName}, nil)
}

func UpdateItemStatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		writeJSON(r.Context(), http.StatusMethodNotAllowed, w, Response{Message: "Method Not Allowed: " + r.Method + ", expecting: " + http.MethodPatch}, nil)
		return
	}
	var ur UpdateItemStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		writeJSON(r.Context(), http.StatusBadRequest, w, Response{Message: "Invalid JSON"}, err)
		return
	}

	response := <-callActor(r, actors.Request{Operation: "UpdateItemStatus", TodoListName: ur.TodoListName, ItemName: ur.ItemName, ItemStatus: ur.ItemStatus})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err updating list"}, response.Err)
		return
	}

	writeJSON(r.Context(), http.StatusCreated, w, Response{Message: "Updated TodoList: " + ur.TodoListName}, nil)
}

func DynamicListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	response := <-callActor(r, actors.Request{Operation: "GetList", TodoListName: r.PathValue("listname")})

	if response.Err != nil {
		writeJSON(r.Context(), http.StatusConflict, w, Response{Message: "Err retrieving list"}, response.Err)
		return
	}

	tmpl := template.Must(template.ParseFiles("./web/dynamic/list.html"))
	tmpl.Execute(w, response.List)
}
