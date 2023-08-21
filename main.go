package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var args struct {
	port  *string
	debug *bool
}

func init() {
	args.port = flag.String("port", "3456", "HTTP server port")
	args.debug = flag.Bool("d", false, "whether or not debug is on")
}

func main() {
	flag.Parse()
	port := os.Getenv("PORT")
	addr := "pass-your-story-on.herokuapp.com"
	prot := "https"
	sprot := "wss"

	if port == "" {
		port = *args.port
		addr = fmt.Sprintf("localhost:%s", port)
		prot = "http"
		sprot = "ws"
	}

	rand.Seed(time.Now().UnixNano())

	kill := make(chan string)
	go KillGameListener(kill)

	hand := &Hand{
		ServerURL:      addr,
		Protocol:       prot,
		Port:           port,
		SocketProtocol: sprot,
		Kill:           kill,
	}

	startServers(hand)
}

func startServers(hand *Hand) {
	h := mux.NewRouter()

	h.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./html/public"))))

	h.HandleFunc("/create", hand.HandleCreate)
	h.HandleFunc("/join", hand.HandleJoin)
	h.HandleFunc("/game", hand.HandleGame)
	h.HandleFunc("/play", hand.HandlePlay)
	h.HandleFunc("/", hand.HandleHomePage)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", hand.Port),
	}
	log.Printf("Starting HTTP server on port %s", hand.Port)
	log.Fatal(srv.ListenAndServe())
}

// Middleware


