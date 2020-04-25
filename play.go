package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"strings"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var Games = make(map[string]*Game, 0)
var Players = make(map[string]*Player, 0)

type Player struct {
	Name       string
	Host       bool
	ToClient   chan string
	Conn       *websocket.Conn
}

/* Play functions */

const NameRequest = `need_name %s`

const NameResponse = `name_input`

func HandlePlayer(gameId string, conn *websocket.Conn) {
	if GameExists(gameId) {
		nameRequest := fmt.Sprintf(NameRequest, "1")
		for {
			err := conn.WriteMessage(1, []byte(nameRequest))
			if err != nil {
				log.Printf("name request write: %s", err)
				return
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("name request read: %s", err)
				return
			}
			nameMessage := string(message)
			split := strings.SplitN(nameMessage, " ", 2)
			if len(split) > 1 && split[0] == NameResponse && len(split[1]) > 0 {
				name := split[1]
				player := FindPlayer(name, gameId)
				if player != nil {
					if player.Conn == nil {
						player.Conn = conn
						break
					}
					nameRequest = fmt.Sprintf(NameRequest, "2")
				} else {
					err = NewPlayer(name, gameId, conn)
					if err == nil {
						break
					} else {
						nameRequest = fmt.Sprintf(NameRequest, "2")
					}
				}
			} else {
				nameRequest = fmt.Sprintf(NameRequest, "0")
			}
		}
	} else {

	}
	//for {
	//	mt, message, err := c.ReadMessage()
	//	if err != nil {
	//		log.Println("read:", err)
	//		break
	//	}
	//	log.Printf("recv: %s", message)
	//	err = c.WriteMessage(mt, message)
	//	if err != nil {
	//		log.Println("write:", err)
	//		break
	//	}
	//}
}

func ExecuteGame(id string) {

}

/* Player Channels */

func (p *Player) ListenFromClient(intoGame chan<- Message) {
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			log.Printf("message read: %s", err)
			break
		}
		m := string(message)
		intoGame <- Message{
			Name: p.Name,
			Content: m,
		}
	}
	p.Conn = nil
}

func (p *Player) SendToClient() {
	for {
		select {
		case m := <-p.ToClient:
			if p.Conn != nil {
				err := p.Conn.WriteMessage(1, []byte(m))
				if err != nil {
					log.Printf("message write: %s", err)
					p.Conn = nil
				}
			}
		}
	}
}

func OutToPlayer(outFromGame <-chan Message, id string) {
	for {
		select {
		case m := <-outFromGame:
			player := FindPlayer(m.Name, id)
			if player != nil {
				player.ToClient <- m.Content
			}
		}
	}
}

/* Maintain Players Array */

func CreatePlayerKey(name, id string) string {
	return fmt.Sprintf("%s%s", name, id)
}

func PlayerExists(name, id string) bool {
	_, ok := Players[CreatePlayerKey(name, id)]
	return ok
}

func FindPlayer(name, id string) *Player {
	if player, ok := Players[CreatePlayerKey(name, id)]; ok {
		return player
	} else {
		return nil
	}
}

func NewPlayer(name, id string, conn *websocket.Conn) error {
	if _, ok := Players[CreatePlayerKey(name, id)]; ok {
		return errors.New("player already exists")
	} else {
		game := FindGame(id)
		host, err := game.AddPlayer(name)
		if err != nil {
			return err
		}
		toClient := make(chan string, 0)
		player := &Player{
			Name: name,
			Host: host,
			ToClient: toClient,
			Conn: conn,
		}
		go player.ListenFromClient(game.Inbound)
		go player.SendToClient()
		Players[CreatePlayerKey(name, id)] = player
		return nil
	}
}

func RemovePlayer(name, id string) {
	if _, ok := Players[CreatePlayerKey(name, id)]; ok {
		delete(Players, CreatePlayerKey(name, id))
	}
}

/* Maintain Games Array */

func GameExists(id string) bool {
	_, ok := Games[id]
	return ok
}

func FindGame(id string) *Game {
	if game, ok := Games[id]; ok {
		return game
	} else {
		return nil
	}
}

func NewGame() string {
	id := CreateID()
	_, ok := Games[id]
	for ok {
		id = CreateID()
		_, ok = Games[id]
	}
	game := CreateGame()
	go OutToPlayer(game.OutBound, id)
	Games[id] = game
	return id
}

func RemoveGame(id string) {
	if _, ok := Games[id]; ok {
		delete(Games, id)
	}
}

func CreateID() string {
	return "PYSO"
	//return fmt.Sprintf("%q%q%q%q", chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))], chars[rand.Intn(len(chars))])
}