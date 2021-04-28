package domain

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMarshalJSON(t *testing.T) {
	c := qt.New(t)

	c.Run("With Regexp", func(c *qt.C) {
		wp, err := NewWebsiteParams("http://foo.org", "", "foo*")
		c.Assert(err, qt.IsNil)
		c.Assert(`{"id": "55065fa3a951948bbb31caf615859b0dbedbb8c5",
                           "url": "http://foo.org",
                           "method": "GET",
                           "match_regexp": "foo*"}`,
			qt.JSONEquals,
			wp)
	})

	c.Run("Without Regexp", func(c *qt.C) {
		wp, err := NewWebsiteParams("http://foo.baz", "HEAD", "")
		c.Assert(err, qt.IsNil)
		c.Assert(`{"id": "4afd9ae157f6fda38a00b0e778e5e98e599ce381",
                           "url": "http://foo.baz",
                           "method": "HEAD",
                           "match_regexp": null}`,
			qt.JSONEquals,
			wp)
	})

}
