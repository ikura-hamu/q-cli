package config

import "github.com/google/uuid"

type Webhook interface {
	GetWebhookID() (string, error)
	GetHostName() (string, error)
	GetSecret() (string, error)
	GetChannels() (map[string]uuid.UUID, error)
}
