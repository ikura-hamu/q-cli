package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/ikura-hamu/q-cli/internal/cmd"
	"github.com/ikura-hamu/q-cli/internal/conf/conf"
)

var version = ""

var option struct {
	Interactive bool
	Version     bool
	Help        bool
	Stdin       bool
	CodeBlock   string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "q [option] [message]")
		flag.PrintDefaults()
	}

	flag.BoolVar(&option.Interactive, "i", false, "Interactive mode")
	flag.BoolVar(&option.Version, "v", false, "Print version")
	flag.BoolVar(&option.Help, "h", false, "Print this message")
	flag.BoolVar(&option.Stdin, "s", false, "Accept input from stdin")
	flag.StringVar(&option.CodeBlock, "c", "", "Send message as code block. Specify the language name (e.g. go, python, shell)")

	flag.Parse()

	if option.Version {
		printVersion()
		return
	}
	if option.Help {
		flag.Usage()
		return
	}

	conf := conf.NewClientConfig()
	client, err := webhook.NewWebhookClient(conf)
	if err != nil {
		log.Fatalf("Failed to create webhook client: %v\n", err)
	}
	command := cmd.NewCommand(client)

	messageOpts := make([]cmd.MessageOptionFunc, 0, 1)
	if option.CodeBlock != "" {
		messageOpts = append(messageOpts, cmd.WithCodeBlock(option.CodeBlock))
	}

	if option.Interactive {
		err = command.SendWithInteractiveMode(messageOpts...)
		if err != nil {
			log.Fatalf("Failed to send message: %v\n", err)
		}
		return
	}
	if option.Stdin {
		err = command.SendWithStdin(messageOpts...)
		if err != nil {
			log.Fatalf("Failed to send message: %v\n", err)
		}
		return
	}

	err = command.SendMessage(strings.Join(flag.Args(), " "), messageOpts...)
	if err != nil {
		log.Fatalf("Failed to send message: %v\n", err)
	}
}

func printVersion() {
	v := ""
	if version != "" {
		v = version
	} else {
		i, ok := debug.ReadBuildInfo()
		if !ok {
			v = "unknown"
		} else {
			v = i.Main.Version
		}
	}
	fmt.Printf("q version %s\n", v)
}
