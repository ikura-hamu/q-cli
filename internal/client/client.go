package client

import (
	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/config"
)

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE}.go . Client

type Client interface {
	// SendMessage sends a message to the webhook URL
	// if message is empty, it should return ErrEmptyMessage
	SendMessage(message string, channelID uuid.UUID) error
}

type Factory[T Client] func(conf config.Webhook) (T, error)
