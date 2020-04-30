package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"net/http"
	"pass-your-story-on/play"
	"strings"
)

func (hand *Hand) DecodeJSONRequest(r *http.Request, required bool) (*PageData, *WebError) {
	decoder := json.NewDecoder(r.Body)
	var req GameRequest
	err := decoder.Decode(&req)
	if err != nil && required {
		return nil, &WebError{http.StatusBadRequest, err.Error(), err}
	}
	return &PageData{
		ServerURL:      hand.ServerURL,
		SocketProtocol: hand.SocketProtocol,
		Protocol:       hand.Protocol,
		GameID:         req.ID,
		PlayerName:     req.Name,
	}, nil
}

func (hand *Hand) DecodeFormRequest(r *http.Request, required bool) (*PageData, *WebError) {
	err := r.ParseForm()
	if err != nil {
		return nil, &WebError{500, err.Error(), err}
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	if required && (id == "" || name == "") {
		return nil, &WebError{http.StatusBadRequest, "missing player name or id", nil}
	}
	req := GameRequest{
		ID: id,
		Name: name,
	}
	return &PageData{
		ServerURL:      hand.ServerURL,
		SocketProtocol: hand.SocketProtocol,
		Protocol:       hand.Protocol,
		GameID:         req.ID,
		PlayerName:     req.Name,
	}, nil
}

func (hand *Hand) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pageData, _ := hand.DecodeJSONRequest(r, false)
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		HandleError(w, "home /", err)
		return
	}
	err = t.Execute(w, pageData)
	if err != nil {
		HandleError(w, "home /", err)
		return
	}
}

func (hand *Hand) HandleCreate(w http.ResponseWriter, r *http.Request) {
	pageData, wErr := hand.DecodeJSONRequest(r, true)
	if wErr != nil {
		HandleError(w, "/create", wErr)
		return
	}
	if pageData.PlayerName == "" {
		HandleError(w, "/create", WebError{http.StatusUnprocessableEntity, "missing player name", nil})
		return
	}
	if *args.debug {

	}
	var gameId string
	if *args.debug {
		gameId = "PYSO"
		if !GameExists(gameId) {
			CreateGame(hand.Kill)
		}
	} else {
		gameId = CreateGame(hand.Kill)
	}
	game := FindGame(gameId)
	if game == nil {
		HandleError(w, "/create", errors.New("could not find game"))
		return
	}
	err := game.AddPlayer(pageData.PlayerName)
	if err != nil {
		HandleError(w, "/create", WebError{http.StatusUnprocessableEntity, err.Error(), err})
		return
	}
	res := new(GameRequest)
	res.ID = gameId
	res.Name = pageData.PlayerName
	resBytes, err := json.Marshal(res)
	if err != nil {
		HandleError(w, "/create", err)
		return
	}
	_, err = w.Write(resBytes)
	if err != nil {
		HandleError(w, "/create", err)
		return
	}
}

// todo: make it so someone can re-join
func (hand *Hand) HandleJoin(w http.ResponseWriter, r *http.Request) {
	pageData, wErr := hand.DecodeJSONRequest(r, true)
	if wErr != nil {
		HandleError(w, "/join", wErr)
		return
	}
	if pageData.PlayerName == "" {
		HandleError(w, "/join", WebError{http.StatusUnprocessableEntity, "missing player name", nil})
		return
	}
	if pageData.GameID == "" {
		HandleError(w, "/join", WebError{http.StatusUnprocessableEntity, "missing game id", nil})
		return
	}
	game := FindGame(pageData.GameID)
	if game == nil {
		HandleError(w, "/join", errors.New("could not find game"))
		return
	}
	err := game.AddPlayer(pageData.PlayerName)
	if err != nil {
		HandleError(w, "/join", WebError{http.StatusUnprocessableEntity, err.Error(), err})
		return
	}
	res := new(GameRequest)
	res.ID = pageData.GameID
	res.Name = pageData.PlayerName
	resBytes, err := json.Marshal(res)
	if err != nil {
		HandleError(w, "/join", err)
		return
	}
	_, err = w.Write(resBytes)
	if err != nil {
		HandleError(w, "/join", err)
		return
	}
}

func (hand *Hand) HandleGame(w http.ResponseWriter, r *http.Request) {
	pageData, wErr := hand.DecodeFormRequest(r, true)
	if wErr != nil {
		HandleError(w, "/game", wErr)
		return
	}
	t, err := template.ParseFiles("html/game.html")
	if err != nil {
		HandleError(w, "/game", err)
		return
	}
	err = t.Execute(w, pageData)
	if err != nil {
		HandleError(w, "/game", err)
		return
	}
}

var upgrader = websocket.Upgrader{}

const Failure = `failed %s`

func (hand *Hand) HandlePlay(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("/play %s", err)
		return
	}
	defer conn.Close()
	err = HandleNewConnection(conn)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(fmt.Sprintf(Failure, err)))
	}
}

const NameSend = `playerName`
const GameIDSend = `gameId`

func HandleNewConnection(conn *websocket.Conn) error {
	name := ""
	gameId := ""
	for name == "" || gameId == "" {
		_, m, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		message := string(m)
		split := strings.SplitN(message, " ", 2)
		if len(split) > 1 && len(split[1]) > 0 {
			switch split[0] {
			case NameSend:
				name = split[1]
			case GameIDSend:
				gameId = split[1]
			}
		}
	}
	game := FindGame(gameId)
	if game == nil {
		return errors.New("game ID not found")
	}
	player := game.FindPlayer(name)
	if player == nil {
		return errors.New("player not found")
	}
	if player.Conn != nil {
		return errors.New("player already exists")
	}
	player.Conn = conn
	player.SendMessage(play.Connected, "")
	if player.Host {
		player.SendMessage(play.Host, "")
	}
	if player.LastControl != "" {
		player.SendMessage(player.LastControl, player.LastBody)
	}
	return player.ReceiveMessages()
}