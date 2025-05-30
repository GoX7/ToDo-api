package mwlogger

import (
	"log/slog"
	"net/http"
	"time"
	"to-do/internal/config"

	"github.com/go-chi/chi/v5/middleware"
)

func New(logger *slog.Logger, cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With(
				slog.String("Method", r.Method),
				slog.String("Path", r.URL.Path),
				slog.String("IP", r.RemoteAddr),
				slog.String("RequestID", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ct := time.Now()

			next.ServeHTTP(ww, r)
			log.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"Request",
				slog.Int("Status", ww.Status()),
				slog.Int("Byte", ww.BytesWritten()),
				slog.String("RequestTime", time.Since(ct).String()),
			)
		}

		return http.HandlerFunc(fn)
	}
}
