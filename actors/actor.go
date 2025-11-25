package actors

import (
	"context"
	"todo/logger"
	"todo/todostore"
)

type Request struct {
	Operation       string
	TodoListName    string
	ItemName        string
	ItemDescription string
	ItemStatus      string
}
type Response struct {
	List todostore.TodoList
	Err  error
}

type Message struct {
	Request      Request
	ResponseChan chan (Response)
	Ctx          context.Context
	Quit         bool
}

var actor *Actor

func GetActor() *Actor {
	if actor == nil {
		actor = &Actor{Name: "Actor1", Messages: make(chan Message, 100)}
	}
	return actor
}

/*********/
/* Actor */
/*********/
type Actor struct {
	Name     string
	Messages chan Message
}

func (a *Actor) SendMessage(ctx context.Context, m Message) {
	logger.InfoLog(ctx, "Actor "+a.Name+" recieved message from trace id: "+m.Ctx.Value(logger.TraceIdKey{}).(string))
	a.Messages <- m
}

func (a *Actor) ProcessMessages(ctx context.Context) {
	for m := range a.Messages {
		logger.InfoLog(ctx, "Actor "+a.Name+" processing message from trace id: "+m.Ctx.Value(logger.TraceIdKey{}).(string)+" operation type: "+m.Request.Operation)

		switch {
		case m.Request.Operation == "CreateList":
			_, err := todostore.CreateList(m.Ctx, m.Request.TodoListName)
			m.ResponseChan <- Response{Err: err}
		case m.Request.Operation == "GetList":
			list, err := todostore.GetList(m.Ctx, m.Request.TodoListName)
			m.ResponseChan <- Response{List: *list, Err: err}
		case m.Request.Operation == "AddItem":
			err := todostore.AddItemToList(m.Ctx, m.Request.TodoListName, m.Request.ItemName, m.Request.ItemDescription)
			m.ResponseChan <- Response{Err: err}
		case m.Request.Operation == "DeleteItem":
			err := todostore.DeleteItemFromList(m.Ctx, m.Request.TodoListName, m.Request.ItemName)
			m.ResponseChan <- Response{Err: err}
		case m.Request.Operation == "UpdateItemDescription":
			err := todostore.UpdateListItemDescription(m.Ctx, m.Request.TodoListName, m.Request.ItemName, m.Request.ItemDescription)
			m.ResponseChan <- Response{Err: err}
		case m.Request.Operation == "UpdateItemStatus":
			err := todostore.UpdateListItemStatus(m.Ctx, m.Request.TodoListName, m.Request.ItemName, m.Request.ItemStatus)
			m.ResponseChan <- Response{Err: err}
		case m.Quit:
			logger.InfoLog(ctx, "Actor "+a.Name+" stopping processing")
			return
		}
	}
}
