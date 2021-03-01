package main

import (
	"github.com/pshvedko/battleship/battle"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"

	"github.com/pshvedko/battleship/api"
	"github.com/pshvedko/battleship/api/ws"
)

func main() {
	b := []byte("TODO_SUPER_SECRET_KEY_1234567890")
	a := api.Application{
		Service: battle.NewBattle(10, 10, 4, 3, 3, 2, 2, 2, 1, 1, 1, 1),
		Logging: log.New(os.Stderr, "", log.LstdFlags),
		Session: sessions.NewCookieStore(b),
	}
	h := mux.NewRouter()
	h.HandleFunc("/begin", a.Begin)
	h.HandleFunc("/click", a.Click)
	h.Use(a.PrepareMiddleware)
	w := ws.WebSocket{
		Updater: websocket.Upgrader{},
		Handler: h,
	}
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
