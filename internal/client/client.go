package client

type Client interface {
	// SendMessage sends a message to the webhook URL
	// if message is empty, it should return ErrEmptyMessage
	SendMessage(message string) error
}
