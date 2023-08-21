package play

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Player struct {
	Name        string          `json:"name"`
	Host        bool            `json:"host"`
	Conn        *websocket.Conn `json:"-"`
	IntoGame    chan<- Message  `json:"-"`
	Story       []Section       `json:"story"`
	OrderIndex  int             `json:"order"`
	LastControl string          `json:"-"`
	LastBody    string          `json:"-"`
	Color       string          `json:"color"`
	TypeFace    string          `json:"typeface"`
}

func (p *Player) ReceiveMessages() error {
	for p.Conn != nil {
		mt, message, err := p.Conn.ReadMessage()
		if err != nil {
			p.Conn = nil
			return err
		}
		if mt == websocket.CloseMessage {
			p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return nil
		}
		m := string(message)
		p.IntoGame <- Message{
			Name:    p.Name,
			Content: m,
		}
	}
	return nil
}

func (p *Player) SendMessage(command, body string) bool {
	switch command {
	case start, write, read, display:
		p.LastControl = command
		p.LastBody = body
	}
	if p.Conn != nil {
		err := p.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s %s", command, body)))
		if err != nil {
			p.Conn = nil
			return false
		}
		return true
	}
	return false
}