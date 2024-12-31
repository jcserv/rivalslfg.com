package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func LogIncomingRequests() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			startTime := time.Now()
			wrapped := wrapResponseWriter(w)

			defer func() {
				if err := recover(); err != nil {
					wrapped.WriteHeader(http.StatusInternalServerError)
					log.Error(ctx, fmt.Sprintf("panic: %v", err))
				}
			}()

			next.ServeHTTP(wrapped, r)

			duration := time.Since(startTime)
			status := wrapped.Status()
			if status == 0 {
				status = 200 // Default to 200 if status wasn't set
			}

			logMsg := fmt.Sprintf("%d | %s %s (%v)",
				status,
				r.Method,
				r.URL.Path,
				duration.Round(time.Millisecond),
			)

			if status >= 500 {
				log.Error(ctx, logMsg)
			} else if status >= 400 {
				log.Warn(ctx, logMsg)
			} else {
				log.Info(ctx, logMsg)
			}
		})
	}
}
