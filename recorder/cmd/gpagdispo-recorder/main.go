package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	env "github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/recorder/pkg/kafka"
)

type config struct {
	KafkaBrokers []string `env:"KAFKA_ADDRS" envDefault:"localhost:9092"`
}

func main() {
	cfg := new(config)

	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("can't parse configuration")
	}

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create Kafka consumer")
	}

	// Gracefully shutdown
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-termChan

		log.Info().Msg("Shutting down...")

		cancel()

		consumer.Close()
	}()

	_ = consumer.Consume(ctx)
}
