package play

import (
	"github.com/pkg/errors"
)

// Adds a player to the game structure
func (g *Game) AddPlayer(name string) error {
	if g.State.Phase == PreGame {
		// check if player exists
		if p, ok := g.Players[name]; ok {
			if p.Conn != nil {
				if p.SendMessage(Alive, "") {
					return errors.New("player already exists")
				}
			}
			// return nil if player lost their connection
			return nil
		}

		if len(g.Order) >= playerLimit {
			return errors.New("player limit reached")
		}

		// Is this the first player created? They are the host.
		host := len(g.Order) == 0
		g.Order = append(g.Order, name)

		player := &Player{
			Name:     name,
			Host:     host,
			IntoGame: g.Messages,
			Story:    make([]Section, 0),
		}

		g.Players[name] = player

		return nil
	} else {
		return errors.New("can't add player after game has started")
	}
}

// Main PreGame Function
func (g *Game) preGameSelect() {
	select {
	case m := <-g.Messages:
		player := g.FindPlayer(m.Name)
		if player == nil {
			return
		}
		command, body := splitMessage(m.Content)
		switch command {
		case StartGame:
			if player.Host {
				g.nextState()
			}
		case Option:
			if player.Host {
				g.updateOption(body)
			}
		default:
			g.commonCommands(player, command, body)
		}
	}
}