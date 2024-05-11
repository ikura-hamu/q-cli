package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type MessageOptionFunc func(string) string

func WithCodeBlock(lang string) MessageOptionFunc {
	return func(message string) string {
		return fmt.Sprintf("```%s\n%s\n```", lang, message)
	}
}

func (cm *Command) SendMessage(message string, options ...MessageOptionFunc) error {
	if message == "" {
		return ErrEmptyMessage
	}

	for _, option := range options {
		message = option(message)
	}

	err := cm.client.SendMessage(message)
	if errors.Is(err, ErrEmptyMessage) {
		return ErrEmptyMessage
	}
	if err != nil {
		return err
	}

	return nil
}

func (cm *Command) SendWithStdin(options ...MessageOptionFunc) error {
	message, err := readFromStdin("")
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	err = cm.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (cm *Command) SendWithInteractiveMode(options ...MessageOptionFunc) error {
	count := 1
	for {
		message, err := readFromStdin(fmt.Sprintf("q [%d]: ", count))
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		if message == "exit" || message == "" {
			fmt.Println()
			break
		}

		for _, option := range options {
			message = option(message)
		}

		err = cm.SendMessage(message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		count++
	}

	return nil
}

func readFromStdin(prompt string) (string, error) {
	var err error
	sb := &strings.Builder{}
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for sc.Scan() {
		fmt.Print(prompt)
		line := sc.Text()
		_, err = sb.WriteString(line + "\n")
		if err != nil {
			return "", fmt.Errorf("failed to write string input: %w", err)
		}
	}

	return sb.String(), nil
}
