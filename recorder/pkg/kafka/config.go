package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
)

// Config defines the Kafka connection configuration
type Config struct {
	TLS struct {
		CAFile   string
		KeyFile  string
		CertFile string
	}
}

// toSaramConfig returns the configuration for Kafka connection
func (kcfg Config) toSaramaConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0

	if kcfg.TLS.CAFile != "" && kcfg.TLS.KeyFile != "" && kcfg.TLS.CertFile != "" {
		keyPair, err := tls.LoadX509KeyPair(kcfg.TLS.CertFile, kcfg.TLS.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("can't load X509 pair: %w", err)
		}
		caCert, err := os.ReadFile(kcfg.TLS.CAFile)
		if err != nil {
			return nil, fmt.Errorf("can't load CA cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			RootCAs:      caCertPool,
		}
	}

	return cfg, nil
}
