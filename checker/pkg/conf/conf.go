// Package conf reads from configuration file to return the list of
// websites to monitor
package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/amzn/ion-go/ion"

	"github.com/sixstone-qq/gpagdispo/checker/pkg/domain"
)

// config defines the format of the configuration file.
type config struct {
	Websites []website `ion:"websites" json:"websites"`
}

// website defines the website params to check in the conf file.
type website struct {
	URL         string `ion:"url" json:"url"`
	Method      string `ion:"method" json:"method"`
	MatchRegexp string `ion:"match_regexp" json:"match_regexp"`
}

// LoadWebsiteParams loads websites to check from a configuration file formatted in ion or JSON.
func LoadWebsiteParams(confPath string) ([]domain.WebsiteParams, error) {
	f, err := os.Open(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var decodeFn func(interface{}) error
	switch path.Ext(confPath) {
	case ".ion":
		decodeFn = ion.NewTextDecoder(f).DecodeTo
	case ".json":
		decodeFn = json.NewDecoder(f).Decode
	default:
		return nil, fmt.Errorf("unknown extension: %s. Supported: ion and json", path.Ext(confPath))
	}

	// Decode
	var cfg config
	if err := decodeFn(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode configuration file: %w", err)
	}

	// Parse content to make sure is correct
	wbParams := make([]domain.WebsiteParams, len(cfg.Websites))
	for i, w := range cfg.Websites {
		params, err := domain.NewWebsiteParams(w.URL, w.Method, w.MatchRegexp)
		if err != nil {
			return nil, fmt.Errorf("can't create website param: %w", err)
		}
		wbParams[i] = *params
	}

	return wbParams, nil
}
