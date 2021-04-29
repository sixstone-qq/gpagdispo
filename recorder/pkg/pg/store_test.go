package pg

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/sixstone-qq/gpagdispo/recorder/pkg/domain"
)

func TestInsertWebsiteResult(t *testing.T) {
	c := qt.New(t)
	ctx := context.Background()

	uri := os.Getenv("POSTGRESQL_DSN")
	if uri == "" {
		uri = "postgres://postgres@localhost/website_test?sslmode=disable"
	}

	s, err := NewStore(uri)
	c.Assert(err, qt.IsNil)
	c.Cleanup(func() { _ = s.Close() })

	err = s.CreateSchema("../../db/migrations")
	c.Assert(err, qt.IsNil)
	c.Cleanup(func() { _ = s.DropSchema("../../db/migrations") })

	// Test data
	wp := domain.WebsiteParams{
		ID:     "id1",
		URL:    "http://foo.org",
		Method: "GET",
	}
	wr := domain.WebsiteResult{
		Elapsed: time.Second,
		Status:  http.StatusOK,
		At:      time.Now().UTC(),
	}

	c.Run("OK", func(c *qt.C) {
		tx, err := s.DB.Beginx()
		c.Assert(err, qt.IsNil)
		defer func() { _ = tx.Rollback() }()

		err = s.InsertWebsiteResult(ctx, wp, wr)
		c.Assert(err, qt.IsNil)

		r := s.DB.QueryRowxContext(ctx,
			`SELECT website_id, elapsed_time, status, matched, unreachable, at
                         FROM websites_results
                         WHERE website_id = $1`, wp.ID)
		c.Assert(r.Err(), qt.IsNil)
		var rec websiteResultRecord
		err = r.StructScan(&rec)
		c.Assert(err, qt.IsNil)
		c.Assert(rec, qt.CmpEquals(cmpopts.EquateApproxTime(time.Second)), websiteResultRecord{
			ID:          "id1",
			Elapsed:     1.0,
			Status:      http.StatusOK,
			Matched:     sql.NullBool{},
			Unreachable: false,
			At:          time.Now().UTC()})
	})

	c.Run("Check no duplicates", func(c *qt.C) {
		tx, err := s.DB.Beginx()
		c.Assert(err, qt.IsNil)
		defer func() { _ = tx.Rollback() }()

		wp.ID = "id2"
		wp.URL = "https://bar.org"
		wp.Method = "HEAD"

		err = s.InsertWebsiteResult(ctx, wp, wr)
		c.Assert(err, qt.IsNil)
		err = s.InsertWebsiteResult(ctx, wp, wr)
		c.Assert(err, qt.IsNil)

		var n int
		err = s.DB.GetContext(ctx, &n,
			`SELECT COUNT(*)
                         FROM websites
                         WHERE id = $1`, wp.ID)
		c.Assert(err, qt.IsNil)
		c.Assert(n, qt.Equals, 1, qt.Commentf("expected inserted websites"))

		err = s.DB.GetContext(ctx, &n,
			`SELECT COUNT(*)
                         FROM websites_results
                         WHERE website_id = $1`, wp.ID)
		c.Assert(err, qt.IsNil)
		c.Assert(n, qt.Equals, 1, qt.Commentf("expected inserted websites results"))
	})
}

type websiteResultRecord struct {
	ID          string       `db:"website_id"`
	Elapsed     float64      `db:"elapsed_time"`
	Status      int          `db:"status"`
	Matched     sql.NullBool `db:"matched"`
	Unreachable bool         `db:"unreachable"`
	At          time.Time    `db:"at"`
}
