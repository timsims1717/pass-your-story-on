package play

import (
	"encoding/json"
	"log"
)

func (g *Game) finishedInit() {
	for _, player := range g.Players {
		player.SendMessage(finish, "")
	}
	g.State.Init = false
	go g.runTimer(g.Options.FinishTimer)
}

func (g *Game) finishedSelect() {
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

// todo: on failure, retry a couple times, then kill the game
func (g *Game) displayInit() {
	stories := make(map[string][]Section, 0)
	for _, player := range g.Players {
		stories[player.Name] = player.Story
	}
	jsonStories, err := json.Marshal(stories)
	if err != nil {
		log.Printf("send display error: %s", err)
		return
	}
	for _, player := range g.Players {
		player.SendMessage(display, string(jsonStories))
	}
	g.State.Init = false
	go g.runTimer(g.Options.KillTimer)
}

func (g *Game) displaySelect() {
	select {
	case <-g.Timing:
		g.Suicide <- g.ID
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