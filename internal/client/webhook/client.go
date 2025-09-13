package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/ras0q/goalie"
)

type WebhookClient struct {
	conf config.Webhook
}

const (
	channelIDHeader string = "X-TRAQ-Channel-ID"
)

func NewWebhookClientFactory(confFactory func() (config.Webhook, error)) func() (*WebhookClient, error) {
	return func() (*WebhookClient, error) {
		conf, err := confFactory()
		if err != nil {
			return nil, fmt.Errorf("create webhook config: %w", err)
		}
		return NewClientFromConfig(conf)
	}
}

func NewClientFromConfig(conf config.Webhook) (*WebhookClient, error) {
	return &WebhookClient{
		conf: conf,
	}, nil
}

func (c *WebhookClient) SendMessage(message string, channelName null.String) (err error) {
	g := goalie.New()
	defer g.Collect(&err)

	if message == "" {
		return client.ErrEmptyMessage
	}

	channelID := uuid.Nil
	if channelName.Valid {
		channels, err := c.conf.GetChannels()
		if err != nil {
			return fmt.Errorf("get channels: %w", err)
		}
		var ok bool
		channelID, ok = channels[channelName.String]
		if !ok {
			return fmt.Errorf("channel '%s' is not found: %w", channelName.String, client.ErrChannelNotFound)
		}
	}

	webhookID, err := c.conf.GetWebhookID()
	if err != nil {
		return fmt.Errorf("get webhook ID: %w", err)
	}

	webhookURLHost, err := c.conf.GetHostName()
	if err != nil {
		return fmt.Errorf("get webhook host: %w", err)
	}
	webhookURL, err := url.JoinPath(webhookURLHost, "/api/v3/webhooks/", webhookID)
	if err != nil {
		return fmt.Errorf("join webhook URL: %w", err)
	}
	webhookURL += "?embed=true"

	req, err := http.NewRequest(http.MethodPost, webhookURL, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	if channelID != uuid.Nil {
		req.Header.Set(channelIDHeader, channelID.String())
	}

	secret, err := c.conf.GetSecret()
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	mac := hmac.New(sha1.New, []byte(secret))

	_, err = mac.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	signature := hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-TRAQ-Signature", signature)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer g.Guard(res.Body.Close)

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message: %s", res.Status)
	}

	return nil
}
