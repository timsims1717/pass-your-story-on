package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(path, name string) (*Page, error) {
	body, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, name))
	if err != nil {
		return nil, WebError{ 500, "could not load page", err }
	}
	return &Page{
		Body: body,
	}, nil
}

type WebError struct {
	Status  int
	Message string
	Cause   error
}

func (e WebError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Status, e.Message)
}