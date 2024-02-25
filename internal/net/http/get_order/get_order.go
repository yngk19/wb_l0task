package get_order

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	response "github.com/yngk19/wb_l0task/internal/utils/api"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	id int `json:"id"`
	response.Response
}

type CacheInterface interface {
	Get(int) interface{}
}

type OrderGetter interface {
	GetOneOrder(context.Context, int) (models.Order, error)
}

func GetOrder(ctx context.Context, log *slog.Logger, orderGettr OrderGetter, cache CacheInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "Time"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			log.Info("order ID is empty")

			render.JSON(w, r, response.Error("not found"))

			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("Invalid order ID: ", err)
		}
		orderFromCache := cache.Get(id)
		if orderFromCache != nil {
			render.JSON(w, r, response.OK(orderFromCache.(models.Order)))
			log.Info(fmt.Sprintf("Order [%s] successfully sent", id))
			return
		}
		order, err := orderGettr.GetOneOrder(ctx, id)
		if err != nil {
			log.Error("Failed to get order", err)
			render.JSON(w, r, response.Error("order does not exist"))
			return
		}
		render.JSON(w, r, response.OK(order))
		log.Info(fmt.Sprintf("Order [%s] successfully sent", id))
		return
	}
}
