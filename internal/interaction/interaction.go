package interaction

import (
	"context"
	"io"
)

type InteractionFactory func() (Interactor, error)

type Interactor interface {
	ReadPassword(prompt string) (line string, err error)
	SetPrompt(prompt string)
	ReadLine() (line string, err error)

	io.Writer
}

type Session interface {
	Run(ctx context.Context, fn func(ctx context.Context, inter Interactor) error) error
}
