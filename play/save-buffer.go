package play

func (g *Game) saveBufferInit() {
	for _, player := range g.Players {
		player.SendMessage(save, "")
	}
	g.State.Init = false
	go g.runTimer(finishTimer)
}

func (g *Game) saveBufferSelect() {
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