package client

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE}.go . Client

type Client interface {
	// SendMessage sends a message to the webhook URL
	// if message is empty, it should return ErrEmptyMessage
	SendMessage(message string) error
}
