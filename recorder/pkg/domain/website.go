package domain

import "time"

// WebsiteParams defines the website parameters to check against
type WebsiteParams struct {
	ID          string  `json:"id" db:"id"`
	URL         string  `json:"url" db:"url"`
	Method      string  `json:"method" db:"method"`
	MatchRegexp *string `json:"match_regexp" db:"match_regexp"`
}

// WebsiteResult defines the result of a website check
type WebsiteResult struct {
	Elapsed time.Duration `json:"elapsed"`
	Status  int           `json:"status"`
	// Matched optionally says if the body response matched the regular expression if provided.
	Matched *bool `json:"matched" `
	// Unreachable means the website check timed out.
	Unreachable bool `json:"unreachable"`
	// At determines when the result was recorded
	At time.Time `json:"at"`
}
