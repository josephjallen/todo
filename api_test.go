package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"todo/actors"
	"todo/logger"
	"todo/todostore"
	"todo/web"

	"github.com/google/uuid"
)

func initActorThread() {
	go func() {
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		logger.InfoLog(ctx, "Starting Actors Thread")
		actors.GetActor().ProcessMessages(ctx)
		logger.InfoLog(ctx, "Actor thread stopped")
	}()
}

/*
go test api_test.go -parallel=2
*/
func TestTodoList(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
	initActorThread()

	t.Run("Create List", func(t *testing.T) {
		logger.InfoLog(nil, "Sending message...")

		mcPostBody := map[string]interface{}{
			"TodoListName": "Shopping",
		}
		body, _ := json.Marshal(mcPostBody)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Go(func() {
				/* CREATE LIST */
				req := httptest.NewRequest(http.MethodGet, "/createlist", bytes.NewReader(body))
				req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
				w := httptest.NewRecorder()
				web.CreateListHandler(w, req)
			})
		}

		wg.Wait()

		/* GET LIST */
		req := httptest.NewRequest(http.MethodGet, "/getlist", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
		w := httptest.NewRecorder()
		web.GetListHandler(w, req)
		var resp todostore.TodoList
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Name != "Shopping" || resp.LItems != nil {
			t.Error("Unexpected list response: " + string(w.Body.Bytes()))
		}
	})

	t.Run("Add To List", func(t *testing.T) {
		logger.InfoLog(nil, "Sending message...")

		mcPostBody := map[string]interface{}{
			"TodoListName": "Shopping",
		}
		body, _ := json.Marshal(mcPostBody)

		/* CREATE LIST */
		req := httptest.NewRequest(http.MethodGet, "/createlist", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
		w := httptest.NewRecorder()
		web.CreateListHandler(w, req)

		mcPostBody = map[string]interface{}{
			"TodoListName":    "Shopping",
			"ItemName":        "Item1",
			"ItemDescription": "Buy Bread",
		}
		body, _ = json.Marshal(mcPostBody)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Go(func() {
				/* ADD ITEM TO LIST */
				req := httptest.NewRequest(http.MethodGet, "/additem", bytes.NewReader(body))
				req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
				w := httptest.NewRecorder()
				web.AddItemHandler(w, req)
			})
		}

		wg.Wait()

		/* GET LIST */
		req = httptest.NewRequest(http.MethodGet, "/getlist", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
		w = httptest.NewRecorder()
		web.GetListHandler(w, req)
		var resp todostore.TodoList
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Name != "Shopping" || len(resp.LItems) != 1 || resp.LItems[0].Name != "Item1" || resp.LItems[0].Description != "Buy Bread" || resp.LItems[0].Status != "not started" {
			t.Error("Unexpected list response: " + string(w.Body.Bytes()))
		}
	})

	t.Cleanup(func() {
		actors.GetActor().Messages <- actors.Message{
			Request: actors.Request{Operation: "quit"},
			Ctx:     ctx,
			Quit:    true,
		}
	})

}
