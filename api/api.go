package api

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/pshvedko/battleship/api/websocket"
	"github.com/pshvedko/battleship/battle"
)

func init() {
	rand.Seed(time.Now().Unix())
	gob.Register(uuid.UUID{})
}

type Application struct {
	Logging *log.Logger
	Session sessions.Store
	Service battle.Battle
}

func (a *Application) SessionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.Session.Get(r, "session")
		for err == nil {
			sid, ok := session.Values["sid"]
			if !ok {
				sid = uuid.New()
				session.Options.SameSite = http.SameSiteLaxMode
				session.Options.HttpOnly = true
				session.Values["sid"] = sid
				err = session.Save(r, w)
				if err != nil {
					break
				}
			}
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "sid", sid)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	})
}

type point struct {
	X int
	Y int
}

type reply struct {
	point
	F int
	C int
}

func (a *Application) Begin(w websocket.ResponseWriter, r *websocket.Request) {
	s := r.Context().Value("sid").(uuid.UUID)
	j := json.NewEncoder(w)
	for _, z := range a.Service.Begin(s) {
		w.WriteHeader(http.StatusOK)
		j.Encode(reply{F: z.F(), point: point{X: z.X(), Y: z.Y()}, C: z.C()})
	}
}

func (a *Application) Click(w websocket.ResponseWriter, r *websocket.Request) {
	s := r.Context().Value("sid").(uuid.UUID)
	var q point
	json.NewDecoder(r.Body).Decode(&q)
	j := json.NewEncoder(w)
	for _, z := range a.Service.Click(s, q.X, q.Y) {
		w.WriteHeader(http.StatusOK)
		j.Encode(reply{F: z.F(), point: point{X: z.X(), Y: z.Y()}, C: z.C()})
	}
}

func (a *Application) Reset(w websocket.ResponseWriter, r *websocket.Request) {
	s := r.Context().Value("sid").(uuid.UUID)
	a.Service.Reset(s)
	w.WriteHeader(http.StatusOK)
}
