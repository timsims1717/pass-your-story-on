package play

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type chat struct {
	Control   bool   `json:"control"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func (g *Game) connectedChat(name string) {
	c := &chat{
		Control: true,
		Source: "",
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05-0700"),
		Message: fmt.Sprintf("%s joined the game", name),
	}
	m, err := json.Marshal(c)
	if err != nil {
		log.Printf("connectedChat error: %s", err)
		return
	}
	for _, player := range g.Players {
		player.SendMessage(Chat, string(m))
	}
}

func (g *Game) chat(name, body string) {
	c := &chat{
		Control: false,
		Source: name,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05-0700"),
		Message: body,
	}
	m, err := json.Marshal(c)
	if err != nil {
		log.Printf("chat error: %s", err)
		return
	}
	for _, player := range g.Players {
		player.SendMessage(Chat, string(m))
	}
}