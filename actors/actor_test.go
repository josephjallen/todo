package actors

import (
	"context"
	"encoding/json"
	"testing"
	"todo/actors"
	"todo/logger"

	"github.com/google/uuid"
)

func initActorThread() {
	go func() {
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		logger.GetCtxLogger(ctx).Info("Starting Actors Thread")
		actors.GetActor().ProcessMessages(ctx)
		logger.GetCtxLogger(ctx).Info("Actor thread stopped")
	}()
}

/*
go test actor_test.go -parallel=2
*/
func TestActors(t *testing.T) {

	initActorThread()
	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	t.Run("Create List", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		request := actors.Request{
			Operation:    "CreateList",
			TodoListName: "Shopping",
		}

		respChannel := make(chan actors.Response)

		actors.GetActor().SendMessage(ctx, actors.Message{
			Request:      request,
			ResponseChan: respChannel,
			Ctx:          ctx,
		})

		response := <-respChannel

		if response.Err != nil && response.Err.Error() != "List already exists: Shopping" {
			t.Error(response.Err.Error())
		}

		request = actors.Request{
			Operation:    "GetList",
			TodoListName: "Shopping",
		}

		actors.GetActor().SendMessage(ctx, actors.Message{
			Request:      request,
			ResponseChan: respChannel,
			Ctx:          ctx,
		})

		response = <-respChannel

		if response.List.Name != "Shopping" {
			t.Error("List not found")
		}
	})

	t.Run("Add To List", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		request := actors.Request{
			Operation:       "AddItem",
			TodoListName:    "Shopping",
			ItemName:        "Item1",
			ItemDescription: "Bread",
		}

		respChannel := make(chan actors.Response)

		actors.GetActor().SendMessage(ctx, actors.Message{
			Request:      request,
			ResponseChan: respChannel,
			Ctx:          ctx,
		})

		response := <-respChannel

		if response.Err != nil && response.Err.Error() != "Item already exists: Item1" {
			t.Error(response.Err.Error())
		}

		request = actors.Request{
			Operation:    "GetList",
			TodoListName: "Shopping",
		}

		actors.GetActor().SendMessage(ctx, actors.Message{
			Request:      request,
			ResponseChan: respChannel,
			Ctx:          ctx,
		})

		response = <-respChannel

		if response.List.Name != "Shopping" {
			t.Error("List not found")
		}

		if len(response.List.LItems) != 1 || response.List.LItems[0].Description != "Bread" {
			resp, _ := json.Marshal(response.List)
			t.Error("Unexpected List Response: " + string(resp))
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
