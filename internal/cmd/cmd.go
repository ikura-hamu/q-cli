package cmd

import "github.com/ikura-hamu/q-cli/internal/client"

type Command struct {
	client client.Client
}

func NewCommand(client client.Client) *Command {
	return &Command{
		client: client,
	}
}
