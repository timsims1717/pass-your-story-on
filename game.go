package main

import (
	"github.com/pkg/errors"
	"math/rand"
)

type GameOptions struct {
	Timer     int
	NumRounds int
	WordLimit int
}

type GameState struct {
	Phase GamePhase
	Round int
}

type Story struct {
	Sections   []Section
	OrderIndex int
}

type Section struct {
	Author  string
	Content string
}

type GamePhase int

const PreGame = GamePhase(0)
const Starting = GamePhase(1)
const ActiveWrite = GamePhase(2)
const ActiveRead = GamePhase(3)
const Finished = GamePhase(4)

type Game struct {
	Players  []string
	Stories  map[string]Story
	Options  *GameOptions
	State    *GameState
	Inbound  chan Message
	OutBound chan Message
}

type Message struct {
	Name string
	Content string
}

func CreateGame() *Game {
	return &Game{
		Options: DefaultOptions(),
		State: &GameState{
			Phase: PreGame,
			Round: 0,
		},
	}
}

/* PreGame */

// Adds a player to the game structure
func (g *Game) AddPlayer(name string) (bool, error) {
	if g.State.Phase == PreGame {
		host := len(g.Players) == 0
		for _, p := range g.Players {
			if p == name {
				return false, errors.New("player already exists")
			}
		}
		g.Players = append(g.Players, name)
		return host, nil
	} else {
		return false, errors.New("can't add player after game has started")
	}
}

// Begins the game.
// Sets the next state, applies the selected options, then randomizes the order of the players.
// Also sets up the Story structures
func (g *Game) StartGame(opt *GameOptions) error {
	if g.State.Phase == PreGame {
		g.NextState()
		g.Options = opt
		players := make([]string, len(g.Players))
		for _, player := range g.Players {
			players = append(players, player)
		}
		g.Players = make([]string, len(g.Players))
		for len(players) > 0 {
			i := rand.Intn(len(players))
			player := players[i]
			g.Players = append(g.Players, player)
			players = append(players[:i], players[i+1:]...)
			g.Stories[player] = Story{
				OrderIndex: i,
			}
		}
		return nil
	} else {
		return errors.New("can't start game after game has started")
	}
}

/* ActiveRead/ActiveWrite */

// Sends out the story so far for each player
func (g *Game) StorySoFar() {
	for _, p := range g.Players {
		//next := g.NextPlayer(p)
		//story := g.Stories[next]
		// convert story to json

		g.OutBound <- Message{
			Name: p,
			Content: "story so far",
		}
	}
}

// Finds the correct author for the current round for a specific player
// Returns "" if the player name is not found
func (g *Game) NextPlayer(player string) string {
	if _, ok := g.Stories[player]; ok {
		index := ((g.State.Round - 1) + g.Stories[player].OrderIndex) % len(g.Players)
		return g.Players[index]
	} else {
		return ""
	}
}

/* All Phases */

// Changes the state of the game based on the current State
func (g *Game) NextState() {
	switch g.State.Phase {
	case PreGame:
		g.State = &GameState{
			Phase: Starting,
			Round: 0,
		}
	case Starting:
		g.State = &GameState{
			Phase: ActiveWrite,
			Round: 1,
		}
	case ActiveWrite:
		round := g.State.Round
		if round == g.Options.NumRounds * len(g.Players) {
			g.State = &GameState{
				Phase: Finished,
				Round: 0,
			}
		} else {
			g.State = &GameState{
				Phase: ActiveRead,
				Round: round,
			}
		}
	case ActiveRead:
		round := g.State.Round
		g.State = &GameState{
			Phase: ActiveWrite,
			Round: round + 1,
		}
	}
}

func (g *Game) SetState(phase GamePhase, round int) {
	g.State = &GameState{
		Phase: phase,
		Round: round,
	}
}

func DefaultOptions() *GameOptions {
	return &GameOptions{
		Timer: 120,
		NumRounds: 1,
		WordLimit: 250,
	}
}