package secret

import (
	"context"
)

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE}.go . SecretDetector

type SecretDetector interface {
	Detect(ctx context.Context, message string) error
}
