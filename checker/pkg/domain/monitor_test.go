package domain

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestMonitor(t *testing.T) {
	c := qt.New(t)

	fetchCounter := int32(0)
	produceCounter := int32(0)
	checker := &Checker{
		FetchWebsiteResult: func(ctx context.Context, wp WebsiteParams) (*WebsiteResult, error) {
			atomic.AddInt32(&fetchCounter, 1)
			return new(WebsiteResult), nil
		},
		ProduceResult: func(wp WebsiteParams, wr WebsiteResult) error {
			atomic.AddInt32(&produceCounter, 1)
			return nil
		},
	}

	wps := []WebsiteParams{{}, {}}

	tick := 100 * time.Millisecond
	timeout := time.Second

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := checker.Monitor(ctx, wps, 100*time.Millisecond)
	c.Assert(err, qt.IsNil)
	c.Assert(int(fetchCounter) <= int(timeout/tick)*len(wps), qt.IsTrue)
	c.Assert(fetchCounter, qt.Equals, produceCounter, qt.Commentf("Same number of fetchs produces same results"))
}
