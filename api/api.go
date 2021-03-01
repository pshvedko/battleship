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
	X int
	Y int
}

type reply struct {
	point
	F int
	C int
}

func (a *Application) Begin(w http.ResponseWriter, r *http.Request) {
	s, ok := r.Context().Value("sid").(uuid.UUID)
	if !ok {
		return
	}
	j := json.NewEncoder(w)
	p := a.Service.Own(s)
	for _, z := range p {
		w.WriteHeader(http.StatusContinue)
		j.Encode(reply{F: 0, point: point{X: z.X(), Y: z.Y()}, C: z.C()})
	}
}

func (a *Application) Click(w http.ResponseWriter, r *http.Request) {
	s, ok := r.Context().Value("sid").(uuid.UUID)
	if !ok {
		return
	}
	var q point
	json.NewDecoder(r.Body).Decode(&q)
	j := json.NewEncoder(w)
	p, c := a.Service.Shot(s, q.X, q.Y)
	for _, z := range p {
		w.WriteHeader(http.StatusContinue)
		j.Encode(reply{F: 0, point: point{X: z.X(), Y: z.Y()}, C: z.C()})
	}
	w.WriteHeader(http.StatusOK)
	j.Encode(reply{F: 1, point: q, C: c})
}
