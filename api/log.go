package api

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

type loggingWriter struct {
	http.ResponseWriter
	header bool
	Status int
	Length int
}

func (w *loggingWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.Length += n
	return
}

func (w *loggingWriter) WriteHeader(code int) {
	if w.header {
		return
	}
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
	w.header = true
	return
}

func (w *loggingWriter) Flush() {
	h, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	h.Flush()
}

func (w *loggingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrHijacked
	}
	return h.Hijack()
}

func (a *Application) LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := &loggingWriter{ResponseWriter: w}
		t := time.Now()
		h.ServeHTTP(l, r)
		if a.Logging != nil {
			a.Logging.Println(r.Method, r.URL, r.Proto, l.Status, l.Length, time.Now().Sub(t))
		}
	})
}
