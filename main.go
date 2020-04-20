package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var args struct {
	port *string
}

func init() {
	args.port = flag.String("port", "3456", "HTTP server port")
}

func main() {
	flag.Parse()
	port := os.Getenv("PORT")
	addr := "pass-your-story-on.herokuapp.com"
	prot := "https"

	if port == "" {
		port = *args.port
		addr = fmt.Sprintf("localhost:%s", port)
		prot = "http"
	}

	startServers(addr, port, prot)
}

func startServers(addr, port, prot string) {
	hand := &Hand{
		ServerURL: addr,
		Protocol:  prot,
	}

	h := mux.NewRouter()

	h.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./html/public"))))

	h.HandleFunc("/", hand.HandleHomePage)
	h.HandleFunc("/create", hand.HandleCreate)
	h.HandleFunc("/game/{id}", hand.HandleGame)
	h.HandleFunc("/play", hand.HandlePlay)

	full := http.TimeoutHandler(installMiddleware(h), 5 * time.Second, "")
	srv := &http.Server {
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: full,
		Addr: fmt.Sprintf(":%s", port),
	}
	log.Printf("Starting HTTP server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

// Middleware

func installMiddleware(h http.Handler) http.Handler {
	return logger(h)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logged := false
		ipList := r.Header.Get("X-Forwarded-For")
		if ipList != "" {
			ips := strings.Split(ipList, ", ")
			if len(ips) > 0 {
				log.Printf("received %s %s request from %s", r.URL, r.Method, ips[0])
				logged = true
			}
		}
		if !logged {
			log.Printf("received %s %s request", r.URL, r.Method)
		}
		h.ServeHTTP(w, r)
	})
}