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
	"github.com/pshvedko/battleship/api/ws"
)

func main() {
	b := []byte("TODO_SUPER_SECRET_KEY_1234567890")
	a := api.Application{
		Logging: log.New(os.Stderr, "", log.LstdFlags),
		Session: sessions.NewCookieStore(b),
		Decoder: schema.NewDecoder(),
	}
	h := mux.NewRouter()
	h.HandleFunc("/begin", a.Begin)
	w := ws.WebSocket{
		Updater: websocket.Upgrader{},
		Handler: h,
	}
	r := mux.NewRouter()
	f := http.FileServer(api.Dir("html"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", f)).Methods(http.MethodGet, http.MethodHead)
	r.Use(w.UpgradeMiddleware)
	r.Use(a.SessionMiddleware)
	r.Use(a.LoggingMiddleware)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
