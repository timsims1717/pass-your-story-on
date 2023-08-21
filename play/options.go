package play

import (
	"encoding/json"
	"log"
)

type GameOptions struct {
	WriteTimer int `json:"writeTimer"`
	ReadTimer  int `json:"readTimer"`
	NumRounds  int `json:"numRounds"`
	WordLimit  int `json:"wordLimit"`
}

func (g *Game) updateOption(options string) {
	var opt *GameOptions
	err := json.Unmarshal([]byte(options), opt)
	if err != nil {
		log.Printf("update options error: %s", err)
		return
	}
	// todo: limit options
	g.Options = opt
	for _, player := range g.Players {
		player.SendMessage(Option, options)
	}
}