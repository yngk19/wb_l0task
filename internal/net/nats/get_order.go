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
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			connectionLostErr = errors.New(fmt.Sprintf("Connection lost, reason: %v", reason))
		}))
	if connectionLostErr != nil {
		return connectionLostErr
	}
	if err != nil {
		connectionLostErr = errors.New(fmt.Sprintf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL))
	}
	log.Info(fmt.Sprintf("Connected to %s clusterID: [%s] clientID: [%s]\n", natsURL, cfg.ClusterID, cfg.ClientID))
	startOpt := stan.StartAt(pb.StartPosition_NewOnly)
	if startSeq != 0 {
		startOpt = stan.StartAtSequence(startSeq)
	} else if deliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if deliverAll && !newOnly {
		startOpt = stan.DeliverAllAvailable()
	} else if startDelta != "" {
		ago, err := time.ParseDuration(startDelta)
		if err != nil {
			sc.Close()
			log.Fatal(err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	subj, i := args[0], 0
	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)
	}

	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj, clientID, qgroup, durable)

	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" || unsubscribe {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
