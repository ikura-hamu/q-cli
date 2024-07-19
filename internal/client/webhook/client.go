package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"net/url"
	"strings"
)

type WebhookClient struct {
	mac        hash.Hash
	client     *http.Client
	webhookURL string
}

func NewWebhookClient(webhookID string, hostName string, secret string) (*WebhookClient, error) {
	mac := hmac.New(sha1.New, []byte(secret))

	client := http.DefaultClient

	webhookURL, err := url.JoinPath(hostName, "/api/v3/webhooks", webhookID)
	if err != nil {
		panic(err)
	}
	webhookURL += "?embed=true"

	return &WebhookClient{
		mac:        mac,
		client:     client,
		webhookURL: webhookURL,
	}, nil
}

func (c *WebhookClient) SendMessage(message string) error {
	req, err := http.NewRequest(http.MethodPost, c.webhookURL, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	c.mac.Reset()
	_, err = c.mac.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	signature := hex.EncodeToString(c.mac.Sum(nil))
	req.Header.Set("X-TRAQ-Signature", signature)

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message: %s", res.Status)
	}

	return nil
}
