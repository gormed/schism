package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/d2r2/go-logger"
	"github.com/gorilla/mux"
)

type TimeOutMiddleware struct {
	logger  logger.PackageLog
	timeout time.Duration
}

func NewTimeOutMiddleware(logger logger.PackageLog, timeout time.Duration) *TimeOutMiddleware {
	return &TimeOutMiddleware{
		logger:  logger,
		timeout: timeout,
	}
}

func (m *TimeOutMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done := make(chan bool)
			ctx, cancelFunc := context.WithTimeout(r.Context(), m.timeout)
			defer cancelFunc()
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()
			select {
			case <-done:
				return
			case <-ctx.Done():
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("request timed out"))
			}
		})
	}
}
