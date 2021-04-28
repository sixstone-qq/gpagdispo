package domain

import "time"

// WebsiteResult defines the result of a website check
type WebsiteResult struct {
	Elapsed time.Duration
	Status  int
	// Matched optionally says if the body response matched the regular expression if provided.
	Matched *bool
	// Unreachable means the website check timed out.
	Unreachable bool
	// At determines when the result was recorded
	At time.Time
}
