package connectors

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// const application_type = "application/json"

type SapphireClient struct {
	client   *http.Client
	endpoint string
	secret   string
}

func NewSapphireClient(endpoint string, secret string) *SapphireClient {
	client := &http.Client{Timeout: 10 * time.Second}
	return &SapphireClient{
		client:   client,
		endpoint: endpoint,
		secret:   secret,
	}
}

func (s *SapphireClient) SendShareToSapphire(body []byte) error {

	req, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("x-shared-secret", s.secret)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		log.Error().Int("status_code", resp.StatusCode).Msg("Non 200 status code")
		return errors.New("return non 200 status code")
	}
	return nil

}
