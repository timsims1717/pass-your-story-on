package main

import (
	"html/template"
	"log"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "testy test",
	}
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		HandleError(w, "home /", err)
		return
	}
	t.Execute(w, page)
}

func HandleError(w http.ResponseWriter, source string, err error) {
	log.Printf("%s %s", source, err)
	if e, ok := err.(WebError); ok {
		http.Error(w, e.Message, e.Status)
	} else {
		http.Error(w, err.Error(), 500)
	}
}