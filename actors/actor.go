package actors

import (
	"fmt"
	"net/http"
)

var actorManager ActorManager = ActorManager{Actors: make(map[string]*Actor)}

type Handler func(http.ResponseWriter, *http.Request)

type Message struct {
	Hand http.Handler
	Resp http.ResponseWriter
	Req  *http.Request
	Chan chan (http.ResponseWriter)
}

type Actor struct {
	Name     string
	Messages []Message
}

func (a *Actor) SendMessage(m Message) {
	a.Messages = append(a.Messages, m)
}

func (a *Actor) ProcessMessages() {
	for _, m := range a.Messages {
		fmt.Printf("Actor %s received message: %s\n", a.Name, m.Req.URL.Path)
		m.Hand.ServeHTTP(m.Resp, m.Req)
		m.Chan <- m.Resp
	}
	a.Messages = nil
}

type ActorManager struct {
	Actors map[string]*Actor
}

func GetActorManager() *ActorManager {
	return &actorManager
}

func (s *ActorManager) RegisterActor(name string) {
	s.Actors[name] = &Actor{Name: name}
}

func (s *ActorManager) SendMessage(m Message) {
	if actor, ok := s.Actors["actor1"]; ok {
		actor.SendMessage(m)
	}
}
