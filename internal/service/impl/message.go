package impl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/ikura-hamu/q-cli/internal/domain/values"
	"github.com/ikura-hamu/q-cli/internal/interaction"
	"github.com/ikura-hamu/q-cli/internal/pkg/types"
	"github.com/ikura-hamu/q-cli/internal/secret"
	"github.com/ikura-hamu/q-cli/internal/service"
)

type Message[Client client.Client] struct {
	clFactory types.Factory[Client]
	rootConf  config.Root
	sec       secret.SecretDetector
	intract   interaction.Session
}

func NewMessage[Client client.Client](cl types.Factory[Client], rootConf config.Root, sec secret.SecretDetector, intract interaction.Session) *Message[Client] {
	return &Message[Client]{
		clFactory: cl,
		rootConf:  rootConf,
		sec:       sec,
		intract:   intract,
	}
}

var _ service.Message = (*Message[client.Client])(nil)

func (m *Message[Client]) Send(ctx context.Context, message values.Message) error {
	cl, err := m.clFactory()
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	if err := m.sec.Detect(ctx, string(message)); err != nil {
		if detectMes, ok := secret.SecretDetected(err); ok {
			fmt.Println(detectMes)
			return nil
		}
		return fmt.Errorf("detect secret: %w", err)
	}

	codeBlock, err := m.rootConf.GetCodeBlock()
	if err != nil {
		return fmt.Errorf("get code block flag: %w", err)
	}

	if codeBlock {
		codeBlockLang, err := m.rootConf.GetCodeBlockLang()
		if err != nil {
			return fmt.Errorf("get code block language: %w", err)
		}

		message = addCodeBlock(message, codeBlockLang)
	}

	printBeforeSend, err := m.rootConf.GetPrintBeforeSend()
	if err != nil {
		return fmt.Errorf("get print before send flag: %w", err)
	}

	if printBeforeSend {
		ok, err := m.checkMessage(ctx, message)
		if err != nil {
			return fmt.Errorf("check message: %w", err)
		}
		if !ok {
			fmt.Println("canceled")
			return nil
		}
	}

	channelName, err := m.rootConf.GetChannelName()
	if err != nil {
		return fmt.Errorf("get channel name: %w", err)
	}

	if err := cl.SendMessage(string(message), channelName); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func addCodeBlock(baseMessage values.Message, codeBlockLang null.String) values.Message {
	leadingBackQuoteCountMax := 0

	for line := range strings.SplitSeq(string(baseMessage), "\n") {
		if !strings.HasPrefix(line, "```") {
			continue
		}
		noLeadingBackQuoteLine := strings.TrimLeft(line, "`")
		leadingBackQuoteCount := len(line) - len(noLeadingBackQuoteLine)
		leadingBackQuoteCountMax = max(leadingBackQuoteCountMax, leadingBackQuoteCount)
	}

	codeBlockBackQuote := strings.Repeat("`", max(leadingBackQuoteCountMax+1, 3))

	return values.Message(fmt.Sprintf("%s%s\n%s\n%s",
		codeBlockBackQuote, codeBlockLang.ValueOrZero(), baseMessage, codeBlockBackQuote))
}

func (m *Message[Client]) checkMessage(ctx context.Context, message values.Message) (bool, error) {
	var ok bool
	err := m.intract.Run(ctx, func(ctx context.Context, inter interaction.Interactor) error {
		fmt.Fprintf(inter, `========Message=========
%s
========================
`, message)

		inter.SetPrompt("Send? [y/n(any)]: ")
		l, err := inter.ReadLine()
		if errors.Is(err, io.EOF) {
			ok = false
			return nil
		}
		if err != nil {
			return fmt.Errorf("interaction read line: %w", err)
		}
		if strings.ToLower(l) != "y" {
			ok = false
			return nil
		}

		ok = true
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("interaction run: %w", err)
	}

	return ok, nil
}
