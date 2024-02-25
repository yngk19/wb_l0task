package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/yngk19/wb_l0task/internal/config"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"log/slog"
)

type Service interface {
	OrderIsValid([]byte, *slog.Logger) (bool, *models.OrderDTO)
}
type Repository interface {
	SaveOrder(context.Context, *models.Order) error
}

type CacheInterface interface {
	Put(int, interface{}) bool
}

func GetOrder(ctx context.Context, cfg *config.Config, log *slog.Logger, orderValidator Service, repo Repository, cache CacheInterface) stan.MsgHandler {
	return func(msg *stan.Msg) {
		log.Info(fmt.Sprintf("NATS: Recieved message [%s] from [%s]", string(msg.Subject), cfg.Nats.ClusterID))
		ok, orderDTO := orderValidator.OrderIsValid(msg.Data, log)
		order := models.Order{ID: 0, Data: *orderDTO}
		if ok {
			err := repo.SaveOrder(ctx, &order)
			if err != nil {
				log.Error("Can't save order: ", err)
				return
			}
			bad := cache.Put(order.ID, order)
			if !bad {
				log.Info(fmt.Sprintf("Message [%d] successfully saved in cache!", order.ID))
			}
			log.Info(fmt.Sprintf("Order [%d] successfully saved!", order.ID))
		}
	}
}
