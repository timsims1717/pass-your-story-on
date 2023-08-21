package play

import (
	"strings"
	"time"
)

func (g *Game) eventLoop() {
	for {
		switch g.State.Phase {
		case PreGame:
			g.preGameSelect()
		case Starting:
			if g.State.Init {
				g.startingInit()
			}
			g.commonSelect()
		case ActiveWrite:
			if g.State.Init {
				g.activeReadWriteInit(false)
			}
			g.activeWriteSelect()
		case SaveBuffer:
			if g.State.Init {
				g.saveBufferInit()
			}
			g.saveBufferSelect()
		case ActiveRead:
			if g.State.Init {
				g.activeReadWriteInit(true)
			}
			g.commonSelect()
		case Finished:
			if g.State.Init {
				g.finishedInit()
			}
			g.finishedSelect()
		case Display:
			if g.State.Init {
				g.displayInit()
			}
			g.displaySelect()
		}
	}
}

func (g *Game) commonSelect() {
	select {
	case <-g.Timing:
		g.nextState()
	case m := <-g.Messages:
		player := g.FindPlayer(m.Name)
		if player == nil {
			return
		}
		command, body := splitMessage(m.Content)
		switch command {
		default:
			g.commonCommands(player, command, body)
		}
	}
}

func (g *Game) commonCommands(player *Player, command, body string) {
	switch command {
	case Connected:
		g.ListPlayersToAll()
		g.connectedChat(player.Name)
	case Chat:
		g.chat(player.Name, body)
	}
}

// Commands
const StartGame = `start`
const Option = `option`
const Chat = `chat`
const Connected = `connected`
const Alive = `alive`
const Host = `host`
const start = `starting`
const write = `write`
const read = `read`
const save = `save`
const finish = `finish`
const display = `display`
const ListPlayers = `listPlayers`
const removePlayer = `removePlayer`
const updatePlayer = `updatePlayer`

// Splits the websocket message into a "command" and a "body"
func splitMessage(message string) (string, string) {
	split := strings.SplitN(message, " ", 2)
	if len(split) > 1 && len(split[1]) > 0 {
		return split[0], split[1]
	}
	return message, ""
}

// Changes the state of the game based on the current State
func (g *Game) nextState() {
	switch g.State.Phase {
	case PreGame:
		g.State = &GameState{
			Phase: Starting,
			Round: 0,
			Init: true,
		}
	case Starting:
		g.State = &GameState{
			Phase: ActiveWrite,
			Round: 1,
			Init: true,
		}
	case ActiveWrite:
		round := g.State.Round
		if round == g.Options.NumRounds*len(g.Order) {
			g.State = &GameState{
				Phase: Finished,
				Round: round,
				Init: true,
			}
		} else {
			g.State = &GameState{
				Phase: SaveBuffer,
				Round: round,
				Init: true,
			}
		}
	case SaveBuffer:
		round := g.State.Round
		g.State = &GameState{
			Phase: ActiveRead,
			Round: round + 1,
			Init: true,
		}
	case ActiveRead:
		round := g.State.Round
		g.State = &GameState{
			Phase: ActiveWrite,
			Round: round,
			Init: true,
		}
	case Finished:
		round := g.State.Round
		g.State = &GameState{
			Phase: Display,
			Round: round,
			Init: true,
		}
	}
}

// Arbitrarily set the state of the game
func (g *Game) setState(phase GamePhase, round int) {
	g.State = &GameState{
		Phase: phase,
		Round: round,
	}
}

func createTimer(add int) int {
	return int(time.Now().Add(time.Duration(add) * time.Second).Unix())
}