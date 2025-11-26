package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
go test api_test.go api.go
*/
func TestTodoList(t *testing.T) {

	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
	initActorThread()
	t.Parallel()

	t.Run("Create List", func(t *testing.T) {
		logger.InfoLog(nil, "Sending message...")

		var result [10]int
		mcPostBody := map[string]interface{}{
			"TodoListName": "Hi",
		}
		body, _ := json.Marshal(mcPostBody)

		for i := 0; i < len(result); i++ {
			go func() {
				j := i
				req := httptest.NewRequest(http.MethodGet, "/createlist", bytes.NewReader(body))
				req = req.WithContext(context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString()))
				w := httptest.NewRecorder()
				web.CreateListHandler(w, req)

				var v web.Response
				json.Unmarshal(w.Body.Bytes(), &v)
				logger.InfoLog(nil, "resp: "+v.Message)

				if v.Message != "Err creating list" && v.Message != "Created List: Hi" {
					result[j] = 3
				} else if todostore.ReadFromMap("Hi") == nil {
					result[j] = 2
				} else {
					result[j] = 1
				}
			}()
		}

		totalSeconds := 0
		for i := true; i == true; {

			routinesNotFinished := 0
			for i := 0; i < len(result); i++ {
				if result[i] == 2 {
					t.Errorf("List not created")
				} else if result[i] == 3 {
					t.Errorf("Unexpected response message")
				} else if result[i] == 0 {
					routinesNotFinished++
				}
			}

			if routinesNotFinished == 0 {
				break
			}

			time.Sleep(1 * time.Second)
			if totalSeconds > 10 {
				t.Errorf("Timeout")
			}
			totalSeconds++
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
