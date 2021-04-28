package conf

import (
	"net/url"
	"os"
	"regexp"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

func TestLoadWebsiteParams(t *testing.T) {
	c := qt.New(t)

	c.Run("OK", func(c *qt.C) {
		c.Run("Empty", func(c *qt.C) {
			cfg, err := LoadWebsiteParams("testdata/empty.json")
			c.Assert(err, qt.IsNil)
			c.Assert(cfg, qt.HasLen, 0)
		})

		expectedWebsiteParams := []domain.WebsiteParams{
			{
				URL:    url.URL{Scheme: "http", Host: "foo.org"},
				Method: domain.HTTPMethodGet,
				ID:     "f068f4ce3120b1e19291215f6e3bab81c6d9aaaf",
			},
			{
				URL:         url.URL{Scheme: "https", Host: "duckduckgo.com", Path: "/search"},
				Method:      domain.HTTPMethodGet,
				MatchRegexp: regexp.MustCompile("duck$"),
				ID:          "fe1a74a16f4978b6b2dea8c3496ae609d3a49ae7",
			},
			{
				URL:         url.URL{Scheme: "http", Host: "only-heads.org", Path: "/foo/bar", RawQuery: "quux=1"},
				Method:      domain.HTTPMethodHead,
				MatchRegexp: regexp.MustCompile("foobar.*"),
				ID:          "32993fbbda453fc52d42b9d74a84d3fe625b6183",
			},
		}

		c.Run("Valid JSON", func(c *qt.C) {
			cfg, err := LoadWebsiteParams("testdata/valid.json")
			c.Assert(err, qt.IsNil)
			c.Assert(cfg, websiteParamsEquals, expectedWebsiteParams)
		})

		c.Run("Valid ION", func(c *qt.C) {
			cfg, err := LoadWebsiteParams("testdata/valid.ion")
			c.Assert(err, qt.IsNil)
			c.Assert(cfg, websiteParamsEquals, expectedWebsiteParams)
		})
	})

	c.Run("NOK", func(c *qt.C) {
		tests := []struct {
			Name      string
			InContent string
			Error     string
		}{
			{
				Name:      "empty params",
				InContent: `{ "websites": [{}] }`,
				Error:     `can't create website param: only HTTP protocol is supported: "" provided`,
			},
			{
				Name:      "wrong scheme",
				InContent: `{ "websites": [{url: "ftp://foo"}] }`,
				Error:     `can't create website param: only HTTP protocol is supported: "ftp" provided`,
			},
			{
				Name:      "wrong method",
				InContent: `{ "websites": [{url: "http://foo.org", method: "POST"}] }`,
				Error:     `can't create website param: can't create HTTP method: unknown HTTP method "POST". Valid ones: \[GET HEAD\]`,
			},
			{
				Name:      "wrong regexp",
				InContent: `{ "websites": [{url: "http://foo.org", method: "HEAD", match_regexp: "["}] }`,
				Error:     `.*error parsing regexp: missing closing ].*`,
			},
		}
		for _, st := range tests {
			c.Run(st.Name, func(c *qt.C) {
				f, err := os.CreateTemp("testdata", "*.ion")
				c.Assert(err, qt.IsNil)
				defer func() {
					f.Close()
					err := os.Remove(f.Name())
					c.Check(err, qt.IsNil)
				}()

				_, err = f.WriteString(st.InContent)
				c.Assert(err, qt.IsNil)

				cfg, err := LoadWebsiteParams(f.Name())
				c.Assert(err, qt.ErrorMatches, st.Error)
				c.Assert(cfg, qt.IsNil)
			})
		}
	})
}

var websiteParamsEquals = qt.CmpEquals(
	cmp.Comparer(func(x, y *regexp.Regexp) bool {
		if x == nil && y == nil {
			return true
		}
		if x == nil && y != nil {
			return false
		}
		if x != nil && y == nil {
			return false
		}
		// Check string representation of Regexp
		return x.String() == y.String()
	}),
)
