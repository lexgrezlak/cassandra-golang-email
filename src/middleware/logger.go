package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("err: %v", err)
				log.Printf("trace: %v", debug.Stack())
			}
		}()

		sw := &statusWriter{ResponseWriter: w}

		start := time.Now()
		// We have to next it before logging, so that we get the status code.
		// Since, the duration numbers are higher.
		next.ServeHTTP(sw, r)
		log.Println("--------------")
		log.Printf("status: %v", sw.status)
		log.Printf("method: %v", r.Method)
		log.Printf("path: %v", r.URL.EscapedPath())
		log.Printf("duration: %v", time.Since(start))
	})
}
