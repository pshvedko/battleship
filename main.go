package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pshvedko/battleship/api"
	"github.com/pshvedko/battleship/api/websocket"
	"github.com/pshvedko/battleship/battle"
)

func main() {
	b := []byte("TODO_SUPER_SECRET_KEY_1234567890")
	a := api.Application{
		Service: battle.NewBattle(4, 3, 3, 2, 2, 2, 1, 1, 1, 1),
		Logging: log.New(os.Stderr, "", log.LstdFlags),
		Session: sessions.NewCookieStore(b),
	}
	w := websocket.New()
	w.HandleFunc("/begin", a.Begin)
	w.HandleFunc("/click", a.Click)
	w.HandleFunc("/reset", a.Reset)
	r := mux.NewRouter()
	f := http.FileServer(api.Dir("html"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", f)).Methods(http.MethodGet, http.MethodHead)
	r.Use(a.LoggingMiddleware)
	r.Use(a.SessionMiddleware)
	r.Use(w.UpgradeMiddleware)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
