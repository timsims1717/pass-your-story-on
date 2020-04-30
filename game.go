package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"pass-your-story-on/play"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var Games = make(map[string]*play.Game, 0)

/* Maintain Games Array */

func GameExists(id string) bool {
	_, ok := Games[id]
	return ok
}

func FindGame(id string) *play.Game {
	if game, ok := Games[id]; ok {
		return game
	} else {
		return nil
	}
}

func CreateGame(kill chan<- string) string {
	id := CreateID()
	_, ok := Games[id]
	for ok {
		id = CreateID()
		_, ok = Games[id]
	}
	game := play.NewGame(id, kill, *args.debug)
	Games[id] = game
	return id
}

func CreateID() string {
	if *args.debug {
		return "PYSO"
	} else {
		return fmt.Sprintf("%q%q%q%q", chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))])
	}
}

func RemoveGame(id string) {
	if _, ok := Games[id]; ok {
		delete(Games, id)
	}
}

func KillGameListener(kill <-chan string) {
	for {
		select {
		case id := <-kill:
			game := FindGame(id)
			if game != nil {
				for _, player := range game.Players {
					if player.Conn != nil {
						player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
						player.Conn = nil
					}
				}
				RemoveGame(id)
				game = nil
			}
		}
	}
}

