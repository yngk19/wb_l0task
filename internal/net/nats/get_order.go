package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/pkg/errors"
	"github.com/yngk19/wb_l0task/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func Connect(cfg config.Nats, log *slog.Logger) error {
	opts := []nats.Option{nats.Name("orders")}
	natsURL := fmt.Sprintf("nats://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return err
	}
	defer nc.Close()
	var connectionLostErr error
	sc, err := stan.Connect(
		cfg.ClusterID, 
		cfg.ClientID,
		stan.NatsConn(nc),
		stan.SetConnectionLostHandler(ConnectionHandler(&connectionLostErr))
	)
	if connectionLostErr != nil {
		return connectionLostErr
	}
	log.Info(fmt.Sprintf("Connected to %s clusterID: [%s] clientID: [%s]\n", natsURL, cfg.ClusterID, cfg.ClientID))
	sub, err := sc.Subscribe("foo", GetOrder(cfg config.Nats), stan.DeliverAllAvailable()
	)
	if err != nil {
		return err
	}
	return nil
}

func ConnectionHandler(connErr *error) stan.ConnectionLostHandler {
	return func(_ stan.Conn, reason error) {
			connErr = &errors.New(fmt.Sprintf("Connection lost, reason: %v", reason))
	}
}

func GetOrder() stan.MsgHandler{
	return func(m *stan.Msg) {
    	log.With(
    		slog.String("Received a message from: ", cfg.ClusterID),
    		slog.String("Subject: ", m.Subject)
    	)
	}
}

