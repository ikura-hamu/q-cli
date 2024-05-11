package conf

import "errors"

var (
	ErrWebhookIDNotSet     = errors.New("webhook ID not set")
	ErrWebhookSecretNotSet = errors.New("webhook secret not set")
)
