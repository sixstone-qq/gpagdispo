package http

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

func TestFetchWebsiteResult(t *testing.T) {
	c := qt.New(t)
	fetcher := &Fetcher{Client: http.DefaultClient}

	c.Run("OK", func(c *qt.C) {
		svr, _ := NewFakeServer(c)
		defer svr.Close()

		c.Run("No regexp", func(c *qt.C) {
			wp, err := domain.NewWebsiteParams(svr.URL, "", "")
			c.Assert(err, qt.IsNil)

			wr, err := fetcher.FetchWebsiteResult(context.TODO(), *wp)
			c.Assert(err, qt.IsNil)
			c.Assert(wr,
				websiteResultEquals,
				&domain.WebsiteResult{At: time.Now().UTC(),
					Status:  http.StatusOK,
					Elapsed: time.Second})
		})

		c.Run("Regexp", func(c *qt.C) {
			wp, err := domain.NewWebsiteParams(svr.URL, http.MethodGet, "Gr")
			c.Assert(err, qt.IsNil)

			wr, err := fetcher.FetchWebsiteResult(context.TODO(), *wp)
			c.Assert(err, qt.IsNil)
			yeah := true
			c.Assert(wr,
				websiteResultEquals,
				&domain.WebsiteResult{At: time.Now().UTC(),
					Status:  http.StatusOK,
					Matched: &yeah,
					Elapsed: time.Second})

			wp, err = domain.NewWebsiteParams(svr.URL, "", "not match")
			c.Assert(err, qt.IsNil)

			wr, err = fetcher.FetchWebsiteResult(context.TODO(), *wp)
			c.Assert(err, qt.IsNil)
			yeah = false
			c.Assert(wr,
				websiteResultEquals,
				&domain.WebsiteResult{At: time.Now().UTC(),
					Status:  http.StatusOK,
					Matched: &yeah,
					Elapsed: time.Second})
		})
	})

	c.Run("Slow", func(c *qt.C) {
		svr, fs := NewFakeServer(c)

		fs.ProcessingTime = time.Minute

		wp, err := domain.NewWebsiteParams(svr.URL, http.MethodGet, "Gr")
		c.Assert(err, qt.IsNil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		c.Cleanup(cancel)

		wr, err := fetcher.FetchWebsiteResult(ctx, *wp)
		c.Assert(err, qt.IsNil)
		c.Assert(wr,
			websiteResultEquals,
			&domain.WebsiteResult{At: time.Now().UTC(),
				Unreachable: true,
				Elapsed:     time.Second})
	})
}

var websiteResultEquals = qt.CmpEquals(
	cmp.Comparer(func(x, y time.Duration) bool {
		return math.Abs(float64(x-y)) < float64(time.Second) && x != 0 && y != 0
	}),
	cmpopts.EquateApproxTime(time.Second),
)

// Fake server

func NewFakeServer(c *qt.C) (*httptest.Server, *fakeServer) {
	fs := new(fakeServer)
	svr := httptest.NewServer(fs)
	c.Cleanup(svr.Close)

	return svr, fs
}

type fakeServer struct {
	ProcessingTime time.Duration
}

func (s *fakeServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	select {
	case <-time.After(s.ProcessingTime):
	case <-req.Context().Done():
	}
	fmt.Fprintln(w, `Great!`)
}
