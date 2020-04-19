package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Hand struct {
	ServerURL string
}

func (hand *Hand) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		ServerURL: hand.ServerURL,
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

func (hand *Hand) HandleGame(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handleGamePost(w, r)
	} else if r.Method == "GET" {
		handleGameGet(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func handleGamePost(w http.ResponseWriter, r *http.Request) {
	request := new(GameRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		HandleError(w, "/game POST", err)
		return
	}
	// todo: create Game ID
	res := new(GameResponse)
	res.ID = "PYSO"
	res.Name = request.Name
	resBytes, err := json.Marshal(res)
	if err != nil {
		HandleError(w, "/game POST", err)
		return
	}
	_, err = w.Write(resBytes)
	if err != nil {
		HandleError(w, "/game POST", err)
		return
	}
}

func handleGameGet(w http.ResponseWriter, r *http.Request) {
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	HandleError(w, "join /game", err)
	//}
}

func (hand *Hand) HandlePlay(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		ServerURL: hand.ServerURL,
	}
	t, err := template.ParseFiles("html/play.html")
	if err != nil {
		HandleError(w, "/play", err)
		return
	}
	err = t.Execute(w, pageData)
	if err != nil {
		HandleError(w, "/play", err)
		return
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