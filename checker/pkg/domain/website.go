package domain

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

// HTTPMethod defines the valid HTTP methods to use in the checker
type HTTPMethod string

const (
	HTTPMethodGet  HTTPMethod = http.MethodGet
	HTTPMethodHead HTTPMethod = http.MethodHead
)

// NewHTTPMethod creates a new HTTPMethod based on a string
func NewHTTPMethod(in string) (HTTPMethod, error) {
	m := HTTPMethod(in)
	switch m {
	case HTTPMethodGet, HTTPMethodHead:
		return m, nil
	}
	return "", fmt.Errorf(`unknown HTTP method "%s". Valid ones: %s`, in, []HTTPMethod{HTTPMethodGet, HTTPMethodHead})
}

// WebsiteParams defines the website parameters to check against
type WebsiteParams struct {
	ID          string         `json:"id"`
	URL         url.URL        `json:"-"`
	Method      HTTPMethod     `json:"method"`
	MatchRegexp *regexp.Regexp `json:"-"`
}

// NewWebsiteParams creates a new WebsiteParmams parsing input strings.
// An empty rawMethod will set Get HTTP method.
// An empty rawRegexp will not generate any regular expression.
func NewWebsiteParams(rawURL, rawMethod, rawRegexp string) (*WebsiteParams, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("can't parse URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf(`only HTTP protocol is supported: "%s" provided`, u.Scheme)
	}

	m := HTTPMethodGet
	if rawMethod != "" {
		m, err = NewHTTPMethod(rawMethod)
		if err != nil {
			return nil, fmt.Errorf("can't create HTTP method: %w", err)
		}
	}

	wp := &WebsiteParams{
		URL:    *u,
		Method: m,
	}

	if rawRegexp != "" {
		wp.MatchRegexp, err = regexp.Compile(rawRegexp)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp: %w", err)
		}
	}

	// Generate the ID based on struct fields
	wp.ID = fmt.Sprintf("%x", sha1.Sum([]byte(wp.URL.String()+string(wp.Method)+rawRegexp)))

	return wp, nil
}

// MarshalJSON provides custom JSON marshalling.
func (wp *WebsiteParams) MarshalJSON() ([]byte, error) {
	var matchRegexp *string
	if wp.MatchRegexp != nil {
		st := wp.MatchRegexp.String()
		matchRegexp = &st
	}

	// Explained at http://choly.ca/post/go-json-marshalling/
	type Alias WebsiteParams
	return json.Marshal(&struct {
		URL         string  `json:"url"`
		MatchRegexp *string `json:"match_regexp"`
		*Alias
	}{
		URL:         wp.URL.String(),
		MatchRegexp: matchRegexp,
		Alias:       (*Alias)(wp),
	})
}
