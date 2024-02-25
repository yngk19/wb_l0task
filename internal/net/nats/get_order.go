package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/yngk19/wb_l0task/internal/config"
	"github.com/yngk19/wb_l0task/internal/repository"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"log/slog"
)

type Service interface {
	OrderIsValid([]byte, *slog.Logger) (bool, *models.Order)
}
type Repository interface {
	SaveOrder(order models.Order) (int64, error)
}

func GetOrder(cfg *config.Config, log *slog.Logger, orderValidator Service, repo repository.Repository) stan.MsgHandler {
	return func(msg *stan.Msg) {
		log.Info(fmt.Sprintf("NATS: Recieved message [%s] from [%s]", string(msg.Subject), cfg.Nats.ClusterID))
		ok, order := orderValidator.OrderIsValid(msg.Data, log)
		if ok {
			id, err := repo.SaveOrder(*order)
			if err != nil {
				log.Error("Can't save order: ", err)
			}
			log.Info(fmt.Sprintf("Order [%s} successfullt saved!", id))
		}
	}
}
