package api

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(uuid.UUID{})
}

type Application struct {
	Logging *log.Logger
	Session *sessions.CookieStore
	Decoder *schema.Decoder
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

func (a *Application) Begin(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusContinue)
	w.WriteHeader(http.StatusContinue)
	w.WriteHeader(http.StatusContinue)
}
