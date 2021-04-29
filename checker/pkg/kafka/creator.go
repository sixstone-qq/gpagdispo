package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

// CreateTopic creates Kafka topic to produces website checks onto if it does not exist.
func CreateTopic(addrs []string) error {
	admin, err := sarama.NewClusterAdmin(addrs, config())
	if err != nil {
		return fmt.Errorf("can't create cluster admin: %w", err)
	}

	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("can't list topics: %w", err)
	}

	_, ok := topics[websiteTopic]
	if !ok {
		err = admin.CreateTopic(websiteTopic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			log.Info().Str("topic", websiteTopic).Msg("Topic created")
		}
		return err
	}

	return nil
}

// config returns the configuration for Kafka connection
func config() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0

	return cfg
}
