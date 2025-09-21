package terminal

import (
	"context"
	"fmt"
	"os"

	"github.com/ikura-hamu/q-cli/internal/interaction"
	"github.com/ras0q/goalie"
	"golang.org/x/term"
)

var _ interaction.Interactor = (*term.Terminal)(nil)

type Session struct{}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) Run(ctx context.Context, fn func(ctx context.Context, inter interaction.Interactor) error) (err error) {
	g := goalie.New()
	defer g.Collect(&err)

	stdinFd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(stdinFd)
	if err != nil {
		return fmt.Errorf("make terminal raw: %w", err)
	}
	defer g.Guard(func() error {
		err := term.Restore(stdinFd, oldState)
		if err != nil {
			return fmt.Errorf("restore terminal state: %w", err)
		}
		return nil
	})

	t := term.NewTerminal(os.Stdin, "")

	return fn(ctx, t)
}
