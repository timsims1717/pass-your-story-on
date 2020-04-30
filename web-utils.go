package main

import (
	"fmt"
	"log"
	"net/http"
)

type Hand struct {
	ServerURL      string
	Protocol       string
	Port           string
	SocketProtocol string
	Kill           chan string
}

type GameRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PageData struct {
	ServerURL      string
	Protocol       string
	SocketProtocol string
	PlayerName     string
	GameID         string
}

type WebError struct {
	Status  int
	Message string
	Cause   error
}

func (e WebError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Status, e.Message)
}


func HandleError(w http.ResponseWriter, source string, err error) {
	log.Printf("%s %s", source, err)
	if e, ok := err.(WebError); ok {
		http.Error(w, e.Message, e.Status)
	} else {
		http.Error(w, err.Error(), 500)
	}
}