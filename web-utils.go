package main

import (
	"fmt"
)

type PageData struct {
	ServerURL  string
	Protocol   string
	PlayerName string
}

type WebError struct {
	Status  int
	Message string
	Cause   error
}

func (e WebError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Status, e.Message)
}