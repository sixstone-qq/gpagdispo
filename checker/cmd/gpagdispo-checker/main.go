package main

import (
	env "github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/conf"
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

	log.Log().Msgf("%+v", websites)
}
