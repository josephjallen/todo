package actors

import (
	"todo/logger"
	"context"
	"net/http"
)

type Message struct {
	Hand http.Handler
	Resp http.ResponseWriter
	Req  *http.Request
	Chan chan (http.ResponseWriter)
	Quit bool
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
	logger.InfoLog(ctx, "Actor " + a.Name + " recieved message: " + m.Req.URL.Path)
	a.Messages <- m
}

func (a *Actor) ProcessMessages(ctx context.Context) {
	for m := range a.Messages {
		if (m.Quit) {
			logger.InfoLog(ctx, "Actor " + a.Name + " stopping processing")
			break
		}
		logger.InfoLog(ctx, "Actor " + a.Name + " processing message: " + m.Req.URL.Path)
		m.Hand.ServeHTTP(m.Resp, m.Req)
		m.Chan <- m.Resp
	}
}

