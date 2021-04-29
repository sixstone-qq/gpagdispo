package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

// Fetcher is in charge of fetch website results from provided websites.
type Fetcher struct {
	Client *http.Client
}

// FetchWebsiteResult fetches the result to monitor from incoming WebsiteParams
func (f *Fetcher) FetchWebsiteResult(ctx context.Context, wp domain.WebsiteParams) (*domain.WebsiteResult, error) {
	req, err := http.NewRequestWithContext(ctx, string(wp.Method), wp.URL.String(), nil)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	resp, err := f.Client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.WebsiteResult{
				Elapsed:     elapsed,
				Unreachable: true,
				At:          time.Now().UTC(),
			}, nil
		}
		return nil, err
	}

	var matched *bool
	if wp.MatchRegexp != nil {
		blob, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		res := wp.MatchRegexp.Match(blob)
		matched = &res
	}

	return &domain.WebsiteResult{
		Status:  &resp.StatusCode,
		Elapsed: elapsed,
		Matched: matched,
		At:      time.Now().UTC(),
	}, nil

}
