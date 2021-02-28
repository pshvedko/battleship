package api

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

func init() {
	gob.Register(uuid.UUID{})
}

type Application struct {
	Logging *log.Logger
	Session *sessions.CookieStore
	Decoder *schema.Decoder
}

type sessionMiddleware struct {
	s sessions.Store
	h http.Handler
}

func (s sessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := s.s.Get(r, "session")
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
		s.h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "sid", sid)))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (a *Application) SessionMiddleware(handler http.Handler) http.Handler {
	return sessionMiddleware{s: a.Session, h: handler}
}

type messageWriter struct {
	*websocket.Conn
	id []byte
	mu chan bool
	bb bytes.Buffer
	mt int
	st int
	hd http.Header
}

func (w *messageWriter) Header() http.Header {
	return w.hd
}

func (w *messageWriter) Write(data []byte) (int, error) {
	return w.bb.Write(data)
}

func (w *messageWriter) WriteHeader(status int) {
	w.st = status
}

func (w messageWriter) Done() {
	m := bytes.Buffer{}
	if !w.Lock() {
		return
	}
	defer w.Unlock()
	_, err := m.Write(w.id)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(&m, "\n%d\n", w.st)
	if err != nil {
		return
	}
	_, err = w.bb.WriteTo(&m)
	if err != nil {
		return
	}
	_ = w.WriteMessage(w.mt, m.Bytes())
}

func (w messageWriter) Unlock() {
	w.mu <- true
}

func (w messageWriter) Lock() bool {
	return <-w.mu
}

type upgradeAndServe struct {
	u websocket.Upgrader
	h http.Handler
}

func (u upgradeAndServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := u.u.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	z := make(chan bool, 1)
	z <- true
	defer close(z)
	defer func() { <-z }()
	var t int
	var b []byte
	var q *http.Request
	for {
		t, b, err = c.ReadMessage()
		if err != nil {
			return
		}
		switch t {
		case websocket.CloseMessage:
			return
		case websocket.TextMessage:
			p := bytes.Split(b, []byte{'\n'})
			var o io.Reader
			var m []byte
			var s []byte
			switch len(p) {
			case 4:
				o = bytes.NewReader(p[3])
				fallthrough
			case 3:
				s = bytes.TrimSpace(p[2])
				m = bytes.TrimSpace(p[1])
				switch string(m) {
				case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
				default:
					// TODO
					continue
				}
			default:
				// TODO
				continue
			}
			q, err = http.NewRequestWithContext(r.Context(), string(m), string(s), o)
			if err != nil {
				// TODO
				continue
			}
			go u.Do(&messageWriter{Conn: c, id: p[0], mu: z, mt: t, bb: bytes.Buffer{}, hd: http.Header{}, st: 200}, q)
		}
	}
}

func (u upgradeAndServe) Do(w *messageWriter, r *http.Request) {
	u.h.ServeHTTP(w, r)
	w.Done()
}

func (a *Application) UpgradeAndServe(updater websocket.Upgrader, handler http.Handler) http.Handler {
	return upgradeAndServe{u: updater, h: handler}
}

func (a *Application) Begin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"a":123, "b":true, "c":"ok"}`))
}
