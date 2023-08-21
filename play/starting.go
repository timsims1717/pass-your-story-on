package play

import "math/rand"

func (g *Game) startingInit() {
	// randomize the order of all the players
	players := make([]string, 0)
	for _, player := range g.Order {
		players = append(players, player)
	}
	g.Order = make([]string, 0)
	for len(players) > 0 {
		i := rand.Intn(len(players))
		player := players[i]
		g.Players[player].OrderIndex = len(g.Order)
		g.Order = append(g.Order, player)
		players = append(players[:i], players[i+1:]...)
	}
	// send the start message
	for _, player := range g.Players {
		player.SendMessage(start, "")
	}
	// start the internal countdown timer
	go g.runTimer(5)
	g.State.Init = false
}