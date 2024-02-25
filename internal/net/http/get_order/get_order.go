package get_order

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	response "github.com/yngk19/wb_l0task/internal/utils/api"
	"log/slog"
	"net/http"
)

type Request struct {
	id int `json:"id" validate:"required"`
}

type Response struct {
	id int `json:"id"`
	response.Response
}

type OrderGetter interface {
	GetOrderById(id int) (*models.Order, error)
}

func GetOrder(log *slog.Logger, orderGettr OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("handler", "Time"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request", err)
			render.JSON(w, r, response.Error("invalid error"))
			return
		}
		id := req.id
		order, err := orderGettr.GetOrderById(id)
		if err != nil {
			log.Error("Failed to get order", err)
			render.JSON(w, r, response.Error("order does not exist"))
			return
		}
		render.JSON(w, r, response.OK(*order))
		log.Info(fmt.Sprintf("Order [%s] successfully sent", id))
		return
	}
}
