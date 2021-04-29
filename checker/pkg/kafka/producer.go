package kafka

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

const websiteTopic = "website.monitor"

// Producer generates the website checks into a Kafka topic
type Producer struct {
	kfkProducer sarama.AsyncProducer
	wg          sync.WaitGroup
}

type payload struct {
	WebsiteParams domain.WebsiteParams `json:"website"`
	WebsiteResult domain.WebsiteResult `json:"result"`
}

// NewProducer handles the creation of the async producers and associated GoRoutines
func NewProducer(addrs []string) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer(addrs, config())
	if err != nil {
		return nil, fmt.Errorf("can't create async producer: %w", err)
	}

	p := &Producer{
		kfkProducer: producer,
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for err := range p.kfkProducer.Errors() {
			log.Error().Err(err).Msg("Error producing result")
		}
	}()

	return p, nil
}

// Produce produces a website check
func (p *Producer) Produce(wp domain.WebsiteParams, wr domain.WebsiteResult) error {
	blob, err := json.Marshal(&payload{
		WebsiteParams: wp,
		WebsiteResult: wr,
	})
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: websiteTopic,
		Key:   sarama.StringEncoder(wp.ID),
		Value: sarama.ByteEncoder(blob),
	}

	p.kfkProducer.Input() <- message

	return nil
}

// Close closes the resources of the producer
func (p *Producer) Close() {
	p.kfkProducer.AsyncClose()
	p.wg.Wait()
}
