package api

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/pshvedko/battleship/battle"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	gob.Register(uuid.UUID{})
}

type Application struct {
	Logging *log.Logger
	Session *sessions.CookieStore
	Service *battle.Battle
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

func (a *Application) PrepareMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid, ok := r.Context().Value("sid").(uuid.UUID)
		if ok {
			err := r.ParseForm()
			if err != nil {
				return
			}
			r.Form.Add("sid", sid.String())
		}
		h.ServeHTTP(w, r)
	})
}

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type reply struct {
	point
	F int `json:"f"`
	C int `json:"c"`
}

func (a *Application) Begin(w http.ResponseWriter, r *http.Request) {
	g := a.Service.Get(r.Context().Value("sid").(uuid.UUID))
	g.Lock()
	defer g.Unlock()
	j := json.NewEncoder(w)
	for x, f := range g.Field0() {
		for y, c := range f {
			w.WriteHeader(http.StatusContinue)
			j.Encode(reply{F: 0, point: point{X: x, Y: y}, C: c})
		}
	}
}

func (a *Application) Click(w http.ResponseWriter, r *http.Request) {
	g := a.Service.Get(r.Context().Value("sid").(uuid.UUID))
	g.Lock()
	defer g.Unlock()
	var p point
	json.NewDecoder(r.Body).Decode(&p)
	json.NewEncoder(w).Encode(reply{F: 1, point: p, C: 0})
}
