/*
Copyright © 2024 ikura-hamu 104292023+ikura-hamu@users.noreply.github.com
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/ikura-hamu/q-cli/internal/message"
	"github.com/ikura-hamu/q-cli/internal/secret"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	ErrChannelNotFound = errors.New("channel name is not in the configuration")
)

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		TraverseChildren: true,
		Use:              "q [message]",
		Example:          "q print(\"Hello, world!\") -c -l py",
		Short:            "traQ Webhook CLI",
		Long: `"q-cli" は、traQにWebhookを使ってメッセージを送信するためのCLIツールです。設定に基づいてWebhookを送信します。
設定は設定ファイルに記述するか、環境変数で指定することができます。
メッセージは標準入力からも受け取ることができます。
設定ファイルの場所は、何も指定しない場合、$HOME/.q-cli.yaml です。`,
	}
	return rootCmd
}

func NewRoot(rootCmd *cobra.Command, fileConf config.File, rootConf config.Root, webhookConf config.Webhook,
	cl *webhook.WebhookClient, mes message.Message, sec secret.SecretDetector) *cobra.Command {

	// rootCmd represents the base command when called without any subcommands
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if v, err := rootConf.GetVersion(); err != nil {
			return fmt.Errorf("failed to get version flag: %w", err)
		} else if v {
			printVersionInfo()
			return nil
		}

		ctx := cmd.Context()

		codeBlock, err := rootConf.GetCodeBlock()
		if err != nil {
			return fmt.Errorf("get code block: %w", err)
		}

		codeBlockLang, err := rootConf.GetCodeBlockLang()
		if err != nil {
			return fmt.Errorf("get code block lang: %w", err)
		}

		messageStr, err := mes.BuildMessage(args, message.Option{
			CodeBlock:     codeBlock,
			CodeBlockLang: codeBlockLang.String,
		})
		if err != nil {
			return fmt.Errorf("failed to build message: %w", err)
		}

		err = sec.Detect(ctx, messageStr)
		if detectMes, ok := secret.SecretDetected(err); ok {
			fmt.Println(detectMes)
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to detect secret: %w", err)
		}

		// TODO: こっちの処理をclientの方でやる
		channelID := uuid.Nil
		channelName, err := rootConf.GetChannelName()
		if err != nil {
			return fmt.Errorf("get channel name: %w", err)
		}
		if channelName.Valid {
			channels, err := webhookConf.GetChannels()
			if err != nil {
				return fmt.Errorf("get channels: %w", err)
			}
			var ok bool
			channelID, ok = channels[channelName.String]
			if !ok {
				return fmt.Errorf("channel '%s' is not found: %w", channelName.String, ErrChannelNotFound)
			}
		}

		printBeforeSend, err := rootConf.GetPrintBeforeSend()
		if err != nil {
			return fmt.Errorf("get print before send: %w", err)
		}
		if printBeforeSend {
			ok, err := checkMessage(messageStr)
			if err != nil {
				return fmt.Errorf("failed to check message: %w", err)
			}
			if !ok {
				fmt.Println("Send canceled.")
				return nil
			}
		}

		err = cl.SendMessage(messageStr, channelID)
		if errors.Is(err, client.ErrEmptyMessage) {
			return errors.New("empty message is not allowed")
		}
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		return nil
	}

	return rootCmd
}

func checkMessage(message string) (bool, error) {
	fmt.Printf(`========Message:========
%s
========================
`, message)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return false, fmt.Errorf("terminal make raw: %w", err)
	}
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		cobra.CheckErr(err)
	}()

	t := term.NewTerminal(os.Stdin, "")
	t.SetPrompt("Send? [y/n(any)]: ")
	l, err := t.ReadLine()
	if errors.Is(err, io.EOF) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("terminal read line: %w", err)
	}
	if strings.ToLower(l) != "y" {
		return false, nil
	}

	return true, nil
}

var (
	version string
)

func printVersionInfo() {
	v := ""
	if version != "" {
		v = version
	} else {
		i, ok := debug.ReadBuildInfo()
		if !ok {
			v = "unknown"
		} else {
			v = i.Main.Version
			if v == "" {
				v = "unknown"
			}
		}
	}
	fmt.Printf("q version %s\n", v)
}
