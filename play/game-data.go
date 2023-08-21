package play

import (
	"encoding/json"
	"log"
	"time"
)

type GameState struct {
	Phase GamePhase
	Round int
	Init  bool
}

type Section struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type Message struct {
	Name    string
	Content string
}

type GamePhase int

const PreGame = GamePhase(0)
const Starting = GamePhase(1)
const ActiveWrite = GamePhase(2)
const SaveBuffer = GamePhase(3)
const ActiveRead = GamePhase(4)
const Finished = GamePhase(5)
const Display = GamePhase(6)

const playerLimit = 12
const finishTimer = 3
const killTimer = 1200

type Game struct {
	ID       string
	Players  map[string]*Player
	Order    []string
	Options  *GameOptions
	State    *GameState
	Messages chan Message
	Timing   chan struct{}
	Suicide  chan<- string
}

func NewGame(id string, kill chan<- string, debug bool) *Game {
	options := new(GameOptions)
	if debug {
		options = debugOptions()
	} else {
		options = defaultOptions()
	}
	game := &Game{
		ID:      id,
		Players: make(map[string]*Player),
		Order:   make([]string, 0),
		Options: options,
		State: &GameState{
			Phase: PreGame,
			Round: 0,
			Init: false,
		},
		Messages: make(chan Message),
		Timing: make(chan struct{}),
		Suicide: kill,
	}
	go game.eventLoop()
	return game
}

func debugOptions() *GameOptions {
	return &GameOptions{
		WriteTimer: 30,
		ReadTimer: 5,
		NumRounds:  1,
		WordLimit:  250,
	}
}

func defaultOptions() *GameOptions {
	return &GameOptions{
		WriteTimer: 120,
		ReadTimer: 24,
		NumRounds:  1,
		WordLimit:  250,
	}
}

func (g *Game) FindPlayer(name string) *Player {
	if player, ok := g.Players[name]; ok {
		return player
	} else {
		return nil
	}
}

func (g *Game) runTimer(t int) {
	<-time.After(time.Duration(t) * time.Second)
	g.Timing <- struct{}{}
}

func (g *Game) ListPlayers() string {
	playerList := make([]*Player, 0)
	for _, name := range g.Order {
		playerList = append(playerList, g.Players[name])
	}
	players, err := json.Marshal(playerList)
	if err != nil {
		log.Printf("list players error: %s", err)
		return ""
	}
	return string(players)
}

func (g *Game) ListPlayersToAll() {
	playerList := g.ListPlayers()
	for _, player := range g.Players {
		player.SendMessage(ListPlayers, playerList)
	}
}