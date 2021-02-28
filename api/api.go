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
	c      *websocket.Conn
	id     []byte
	mu     chan bool
	buffer bytes.Buffer
	proto  int
	status int
	header http.Header
	dirty  bool
}

func (w *messageWriter) Header() http.Header {
	return w.header
}

func (w *messageWriter) Write(data []byte) (int, error) {
	w.dirty = true
	return w.buffer.Write(data)
}

func (w *messageWriter) WriteHeader(status int) {
	w.Flush()
	w.dirty = true
	w.status = status
}

func (w *messageWriter) Flush() {
	if !w.dirty || !<-w.mu {
		return
	}
	defer func() { w.mu <- true }()
	m := bytes.Buffer{}
	_, err := m.Write(w.id)
	if err == nil {
		_, err = fmt.Fprintf(&m, "\n%d\n", w.status)
		if err == nil {
			_, err = w.buffer.WriteTo(&m)
			if err == nil {
				err = w.c.WriteMessage(w.proto, m.Bytes())
				if err == nil {
					w.dirty = false
					return
				}
			}
		}
	}
	_ = w.c.Close()
}

func (w *messageWriter) Status100() {
	if w.status > 99 && w.status < 200 {
		w.status += 100
	}
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
					return
				}
			default:
				return
			}
			q, err = http.NewRequestWithContext(r.Context(), string(m), string(s), o)
			if err != nil {
				return
			}
			go u.Do(&messageWriter{
				c:      c,
				id:     p[0],
				mu:     z,
				proto:  t,
				buffer: bytes.Buffer{},
				status: http.StatusOK,
				header: http.Header{},
				dirty:  false,
			}, q)
		}
	}
}

func (u upgradeAndServe) Do(w *messageWriter, r *http.Request) {
	u.h.ServeHTTP(w, r)
	w.Status100()
	w.Flush()
}

func (a *Application) UpgradeAndServe(updater websocket.Upgrader, handler http.Handler) http.Handler {
	return upgradeAndServe{u: updater, h: handler}
}

func (a *Application) Begin(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusContinue)
	w.WriteHeader(http.StatusContinue)
	w.WriteHeader(http.StatusContinue)
}
