package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

const topic = "website.monitor"

// Consumer consumes website checks from a Kafka topic
type Consumer struct {
	kfkConsumerGroup sarama.ConsumerGroup
	handler          sarama.ConsumerGroupHandler
	wg               sync.WaitGroup
}

// NewConsumer creates the consumer group from the given addresses
func NewConsumer(addrs []string) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(addrs, "website-monitor-1", cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create consumer: %w", err)
	}

	c := &Consumer{
		kfkConsumerGroup: consumer,
		handler:          new(handler),
	}

	// Track errors
	c.wg.Add(1)
	go func() {
		for err := range c.kfkConsumerGroup.Errors() {
			if errors.Is(err, context.DeadlineExceeded) {
				// Ignore the error
				continue
			}

			log.Error().Err(err).Msg("error consuming result")
		}
	}()

	return c, nil
}

// Consume will loop forever consuming Kafka topics
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		// This method should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims.
		err := c.kfkConsumerGroup.Consume(ctx, []string{topic}, c.handler)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}

// Close closes the resources of the consumer
func (c *Consumer) Close() error {
	err := c.kfkConsumerGroup.Close()
	c.wg.Wait()
	return err
}

// Compile-time check around interface implementatoin
var _ sarama.ConsumerGroupHandler = (*handler)(nil)

// handler implements sarama.ConsumerGroupHandler interface to perform the actions
type handler struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *handler) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (h *handler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (h *handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
loop:
	for {
		select {
		case msg := <-claim.Messages():
			log.Log().Int32("partition", msg.Partition).Int64("offset", msg.Offset).
				Str("key", string(msg.Key)).Str("body", string(msg.Value)).Msg("consumed message")
			sess.MarkMessage(msg, "")
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				log.Info().Err(err).Msg("work done")
			}
			break loop
		}
	}

	return nil
}
