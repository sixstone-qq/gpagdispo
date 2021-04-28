package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/conf"
	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
	chttp "github.com/sixstone-qq/gpagdispo/checker/pkg/http"
	"github.com/sixstone-qq/gpagdispo/checker/pkg/kafka"
)

type config struct {
	ConfigFilePath string   `env:"CONFIG_PATH" envDefault:"websites.ion"`
	KafkaBrokers   []string `env:"KAFKA_ADDRS" envDefault:"localhost:9092"`
}

func main() {
	cfg := new(config)

	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("can't parse configuration")
	}

	websites, err := conf.LoadWebsiteParams(cfg.ConfigFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("can't load file")
	}

	// TODO: Set best client params
	fetcher := &chttp.Fetcher{Client: http.DefaultClient}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create Kafka producer")
	}

	checker := &domain.Checker{
		FetchWebsiteResult: fetcher.FetchWebsiteResult,
		ProduceResult:      producer.Produce,
	}

	// Gracefully shutdown
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-termChan

		log.Info().Msg("Shutting down...")

		cancel()
	}()

	_ = checker.Monitor(ctx, websites, 2*time.Second)
}
