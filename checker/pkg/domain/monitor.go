package domain

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// delayTicks definms the stacked work a worker can have due to delays
const delayTicks = 2

// Checker is in charge of checks website availability
type Checker struct {
	FetchWebsiteResult func(ctx context.Context, wp WebsiteParams) (*WebsiteResult, error)
}

// Monitor periodically checks websites indefinitely
func (c *Checker) Monitor(ctx context.Context, websites []WebsiteParams, tick time.Duration) error {
	// Create the same number of workers than websites
	// This could be optimised to do a pool of workers.
	work := make(map[int]chan WebsiteParams, delayTicks)
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for i := 0; i < len(websites); i++ {
		work[i] = make(chan WebsiteParams)
		go c.worker(i, work[i], &wg, tick)
	}

loop:
	for {
		select {
		case <-time.After(tick):
			for i, wp := range websites {
				work[i] <- wp
			}
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				log.Info().Err(err).Msg("Work done")
			}
			break loop
		}
	}

	// Closing the business
	for i := 0; i < len(websites); i++ {
		close(work[i])
	}

	wg.Wait()

	return nil
}

// worker does the work of perform the HTTP request and produce the events as goroutine.
func (c *Checker) worker(id int, work <-chan WebsiteParams, wg *sync.WaitGroup, maxProcessingTime time.Duration) {
	defer wg.Done()

	for wp := range work {
		log.Log().Int("id", id).Msgf("%+v", wp)
		ctx, cancel := context.WithTimeout(context.Background(), maxProcessingTime)
		wr, err := c.FetchWebsiteResult(ctx, wp)
		cancel()
		log.Log().Err(err).Msgf("Result: %+v", wr)
	}
}
