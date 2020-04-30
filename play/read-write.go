package play

import (
	"encoding/json"
	"log"
)

type StorySend struct {
	Timer int       `json:"timer"`
	Story []Section `json:"story"`
}

func (g *Game) activeReadWriteInit(isRead bool) {
	var timer int
	var command string
	if isRead {
		roundMult := g.State.Round - 1
		if roundMult < 1 {
			roundMult = 1
		}
		timer = g.Options.ReadTimer * roundMult
		command = read
	} else {
		timer = g.Options.WriteTimer
		command = write
	}
	for _, player := range g.Players {
		next := g.Players[g.nextPlayer(player.Name)]
		if next != nil {
			send := StorySend{
				Timer: createTimer(timer),
				Story: next.Story,
			}
			m, err := json.Marshal(send)
			if err != nil {
				log.Printf("send writing error: %s", err)
				continue
			}
			player.SendMessage(command, string(m))
		}
	}
	go g.runTimer(timer + 2)
	g.State.Init = false
}

func (g *Game) activeWriteSelect() {
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
		case save:
			g.saveStory(player, body)
		default:
			g.commonCommands(player, command, body)
		}
	}
}

// Finds the correct author for the current round for a specific player
// Returns "" if the player name is not found
func (g *Game) nextPlayer(player string) string {
	if player, ok := g.Players[player]; ok {
		index := ((g.State.Round - 1) + player.OrderIndex) % len(g.Order)
		return g.Order[index]
	} else {
		return ""
	}
}


func (g *Game) saveStory(player *Player, story string) {
	next := g.Players[g.nextPlayer(player.Name)]
	i := g.State.Round - 1
	if next != nil {
		story := Section{
			Author: player.Name,
			Content: story,
		}
		if len(next.Story) <= i {
			next.Story = append(next.Story, story)
		} else {
			next.Story[i] = story
		}
	}
}