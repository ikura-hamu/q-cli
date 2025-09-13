package client

import (
	"errors"
)

var (
	ErrEmptyMessage    = errors.New("empty message")
	ErrChannelNotFound = errors.New("channel not found")
)
