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
	checker := &Checker{
		FetchWebsiteResult: func(ctx context.Context, wp WebsiteParams) (*WebsiteResult, error) {
			atomic.AddInt32(&fetchCounter, 1)
			return nil, nil
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
}
