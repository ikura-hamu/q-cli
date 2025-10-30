package service

import (
	"context"

	"github.com/ikura-hamu/q-cli/internal/domain/values"
)

type Message interface {
	Send(ctx context.Context, message values.Message) error
}
