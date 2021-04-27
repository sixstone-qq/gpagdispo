package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/conf"
	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

type config struct {
	ConfigFilePath string `env:"CONFIG_PATH" envDefault:"websites.ion"`
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

	checker := &domain.Checker{
		FetchWebsiteResult: func(ctx context.Context, wp domain.WebsiteParams) (*domain.WebsiteResult, error) {
			// TOOD
			return nil, nil
		},
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