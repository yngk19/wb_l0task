package nats

import (
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
