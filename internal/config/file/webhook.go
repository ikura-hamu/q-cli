package file

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/viper"
)

const (
	configKeyWebhookHost   = "webhook_host"
	configKeyWebhookID     = "webhook_id"
	configKeyWebhookSecret = "webhook_secret"
	configKeyChannels      = "channels"
)

type Webhook struct {
	v *viper.Viper
}

var _ config.Webhook = (*Webhook)(nil)

func NewWebhook(v *viper.Viper) *Webhook {
	return &Webhook{
		v: v,
	}
}

func (w *Webhook) GetWebhookID() (string, error) {
	v := w.v.GetString(configKeyWebhookID)
	if v == "" {
		return "", fmt.Errorf("webhook ID is not set")
	}
	return v, nil
}

func (w *Webhook) GetHostName() (string, error) {
	v := w.v.GetString(configKeyWebhookHost)
	if v == "" {
		return "", fmt.Errorf("webhook host is not set")
	}
	return v, nil
}

func (w *Webhook) GetSecret() (string, error) {
	v := w.v.GetString(configKeyWebhookSecret)
	if v == "" {
		return "", fmt.Errorf("webhook secret is not set")
	}
	return v, nil
}

func (w *Webhook) GetChannels() (map[string]uuid.UUID, error) {
	v := w.v.GetStringMapString(configKeyChannels)
	if len(v) == 0 {
		return nil, fmt.Errorf("webhook channels are not set")
	}

	channels := make(map[string]uuid.UUID, len(v))
	for name, idStr := range v {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid channel ID for channel '%s': %w", name, err)
		}
		channels[name] = id
	}

	return channels, nil
}
