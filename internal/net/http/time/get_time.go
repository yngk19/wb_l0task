package get_time

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

func Time(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const handler = "handlers.Time"

		log = log.With(
			slog.String("handler", handler),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		t := time.Now()

		render.JSON(w, r, t)

		log.Info("Time successully sent", slog.String("time", t.String()))

		return
	}
}
