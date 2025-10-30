package client

import (
	"github.com/guregu/null/v6"
)

type Client interface {
	// SendMessage sends a message to the webhook URL
	// if message is empty, it should return ErrEmptyMessage
	SendMessage(message string, channelName null.String) error
}
