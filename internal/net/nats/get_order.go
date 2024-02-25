package nats

import (
	"fmt"
	_ "fmt"
	"github.com/nats-io/stan.go"
	_ "github.com/nats-io/stan.go/pb"
	"github.com/yngk19/wb_l0task/internal/config"
	"log/slog"
	"time"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25
)

func NewNatsConnect(cfg *config.Config, log *slog.Logger) (stan.Conn, error) {
	return stan.Connect(
		cfg.Nats.ClusterID,
		cfg.Nats.ClientID,
		stan.ConnectWait(connectWait),
		stan.PubAckWait(pubAckWait),
		stan.NatsURL("nats://nats:4222"),
		stan.Pings(interval, maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Error("Connection lost, reason: %v", reason)
		}),
		stan.MaxPubAcksInflight(maxPubAcksInflight),
	)
}

func Listen(cfg config.Nats, log *slog.Logger) (stan.Conn, stan.Subscription, error) {
	sc, err := stan.Connect(
		cfg.ClusterID,
		cfg.ClientID,
		stan.NatsURL("nats://admin:123@nats:4222"),
		stan.SetConnectionLostHandler(
			func(_ stan.Conn, reason error) {
				log.Error("Connection with NATS-streaming is lost!")
			},
		),
	)
	if err != nil {
		return nil, nil, err
	}
	connectionInfo := fmt.Sprintf("Connected to nats-streaming server: [%s] clusterID: [%s] clientID]\n", cfg.ClusterID, cfg.ClientID)
	log.Info(connectionInfo)
	mcb := func(msg *stan.Msg) {
		GetOrder(cfg, log, msg)
	}
	sub, err := sc.Subscribe("orders", mcb)
	if err != nil {
		sc.Close()
		return nil, nil, err
	}
	return sc, sub, nil
}

func GetOrder(cfg config.Nats, log *slog.Logger, m *stan.Msg) {
	log.With(
		slog.String("Received a message from: ", cfg.ClusterID),
		slog.String("Subject: ", m.Subject),
	)
}
