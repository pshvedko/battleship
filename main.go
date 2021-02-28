package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"

	"github.com/pshvedko/battleship/api"
)

func main() {
	b := []byte("TODO_SUPER_SECRET_KEY_1234567890")
	a := api.Application{
		Logging: log.New(os.Stderr, "", log.LstdFlags),
		Session: sessions.NewCookieStore(b),
		Decoder: schema.NewDecoder(),
	}
	r := mux.NewRouter()
	q := r.PathPrefix("/api/v1/").Subrouter()
	u := websocket.Upgrader{}
	h := mux.NewRouter()
	h.HandleFunc("/begin", a.Begin)
	q.Handle("/websocket", a.UpgradeAndServe(u, h)).Methods(http.MethodGet)
	f := http.FileServer(api.Dir("html"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", f)).Methods(http.MethodGet, http.MethodHead)
	r.Use(a.SessionMiddleware)
	r.Use(a.LoggingMiddleware)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
