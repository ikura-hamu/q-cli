package conf

type ClientConfig interface {
	GetWebhookID() (string, error)
	GetWebhookSecret() (string, error)
}
