package actors

import (
	"fmt"
	"net/http"
)

var actorManager ActorManager = ActorManager{Actors: make(map[string]*Actor)}

type Message struct {
	Hand http.Handler
	Resp http.ResponseWriter
	Req  *http.Request
	Chan chan (http.ResponseWriter)
}

/*********/
/* Actor */
/*********/
type Actor struct {
	Name     string
	Messages chan Message
}

func (a *Actor) SendMessage(m Message) {
	a.Messages <- m
	fmt.Printf("Actor %s recieved message: %s\n", a.Name, m.Req.URL.Path)
}

func (a *Actor) ProcessMessages() {
	for m := range a.Messages {
		fmt.Printf("Actor %s processing message: %s\n", a.Name, m.Req.URL.Path)
		m.Hand.ServeHTTP(m.Resp, m.Req)
		m.Chan <- m.Resp
	}
}

/*****************/
/* Actor Manager */
/*****************/
type ActorManager struct {
	Actors map[string]*Actor
}

func GetActorManager() *ActorManager {
	return &actorManager
}

func (s *ActorManager) RegisterActor(name string) {
	s.Actors[name] = &Actor{Name: name, Messages: make(chan Message, 100)}
}

func (s *ActorManager) SendMessage(m Message) {
	if actor, ok := s.Actors["actor1"]; ok {
		actor.SendMessage(m)
	}

}
