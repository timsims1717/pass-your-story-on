package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

type Hand struct {
	ServerURL string
	Protocol  string
}

func (hand *Hand) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		ServerURL: hand.ServerURL,
		Protocol:  hand.Protocol,
	}
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

type GameRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GameResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (hand *Hand) HandleCreate(w http.ResponseWriter, r *http.Request) {
	request := new(GameRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		HandleError(w, "/create", err)
		return
	}
	// todo: create Game ID
	res := new(GameResponse)
	res.ID = "PYSO"
	res.Name = request.Name
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

func handleJoin(w http.ResponseWriter, r *http.Request) {
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	HandleError(w, "join /game", err)
	//}
}

func (hand *Hand) HandleGame(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		ServerURL: hand.ServerURL,
		Protocol:  hand.Protocol,
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

func (hand *Hand) HandlePlay(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("%s %s", "/play", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func HandleError(w http.ResponseWriter, source string, err error) {
	log.Printf("%s %s", source, err)
	if e, ok := err.(WebError); ok {
		http.Error(w, e.Message, e.Status)
	} else {
		http.Error(w, err.Error(), 500)
	}
}