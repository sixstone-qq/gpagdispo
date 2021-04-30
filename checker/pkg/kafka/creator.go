package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

// CreateTopic creates Kafka topic to produces website checks onto if it does not exist.
func CreateTopic(addrs []string, cfg Config) error {
	saramaCfg, err := cfg.toSaramaConfig()
	if err != nil {
		return fmt.Errorf("can't create config: %w", err)
	}

	admin, err := sarama.NewClusterAdmin(addrs, saramaCfg)
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
