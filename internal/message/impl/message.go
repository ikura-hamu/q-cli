package impl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ikura-hamu/q-cli/internal/message"
)

type Message struct{}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) BuildMessage(args []string, option message.Option) (string, error) {
	var message string
	var err error

	if len(args) > 0 {
		message = strings.Join(args, " ")
	} else {
		message, err = scan()
		if err != nil {
			return "", fmt.Errorf("failed to scan message: %w", err)
		}
	}

	if option.CodeBlock {
		message = addCodeBlock(message, option.CodeBlockLang)
	}

	return message, nil
}

func scan() (string, error) {
	sc := bufio.NewScanner(os.Stdin)
	sb := &strings.Builder{}
	for sc.Scan() {
		text := sc.Text()
		sb.WriteString(text + "\n")
	}

	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	return strings.TrimSpace(sb.String()), nil
}

func addCodeBlock(baseMessage string, codeBlockLang string) string {
	leadingBackQuoteCountMax := 0

	for _, line := range strings.Split(baseMessage, "\n") {
		if !strings.HasPrefix(line, "```") {
			continue
		}
		noLeadingBackQuoteLine := strings.TrimLeft(line, "`")
		leadingBackQuoteCount := len(line) - len(noLeadingBackQuoteLine)
		leadingBackQuoteCountMax = max(leadingBackQuoteCountMax, leadingBackQuoteCount)
	}

	codeBlockBackQuote := strings.Repeat("`", max(leadingBackQuoteCountMax+1, 3))

	return fmt.Sprintf("%s%s\n%s\n%s", codeBlockBackQuote, codeBlockLang, baseMessage, codeBlockBackQuote)
}
