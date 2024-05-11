package conf

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ikura-hamu/q-cli/internal/conf"
)

type ClientConfig struct {
	configData *configData
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		configData: readFromFile(),
	}
}

type configData struct {
	WebhookID     string `json:"webhook_id"`
	WebhookSecret string `json:"webhook_secret"`
}

func readFromFile() *configData {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	qConfDir := path.Join(userConfDir, "q-cli")
	if _, err := os.Stat(qConfDir); os.IsNotExist(err) {
		return &configData{}
	} else if err != nil {
		panic(err)
	}

	qConfFileName := path.Join(qConfDir, "config.json")
	if _, err := os.Stat(qConfFileName); os.IsNotExist(err) {
		return &configData{}
	} else if err != nil {
		panic(err)
	}

	file, err := os.Open(qConfFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var configData configData
	err = json.NewDecoder(file).Decode(&configData)
	if err != nil {
		panic(err)
	}

	return &configData
}

func (c *ClientConfig) GetWebhookID() (string, error) {
	if c.configData.WebhookID != "" {
		return c.configData.WebhookID, nil
	}

	webhookID, ok := os.LookupEnv("Q_WEBHOOK_ID")
	if !ok {
		return "", conf.ErrWebhookIDNotSet
	}

	return webhookID, nil
}

func (c *ClientConfig) GetWebhookSecret() (string, error) {
	if c.configData.WebhookSecret != "" {
		return c.configData.WebhookSecret, nil
	}

	webhookSecret, ok := os.LookupEnv("Q_WEBHOOK_SECRET")
	if !ok {
		return "", conf.ErrWebhookSecretNotSet
	}

	return webhookSecret, nil
}
