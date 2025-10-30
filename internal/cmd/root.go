/*
Copyright © 2024 ikura-hamu 104292023+ikura-hamu@users.noreply.github.com
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/ikura-hamu/q-cli/internal/domain/values"
	"github.com/ikura-hamu/q-cli/internal/service"
	"github.com/spf13/cobra"
)

var (
	ErrChannelNotFound = errors.New("channel name is not in the configuration")
)

type RootBare struct {
	*cobra.Command
}

type Root struct {
	*cobra.Command
}

func NewRootBare() *RootBare {
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
	return &RootBare{
		Command: rootCmd,
	}
}

func NewRoot[Client client.Client](rootCmd *RootBare, rootConf config.Root, mes service.Message) *Root {

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if v, err := rootConf.GetVersion(); err != nil {
			return fmt.Errorf("get version flag: %w", err)
		} else if v {
			printVersionInfo()
			return nil
		}

		ctx := cmd.Context()

		var message values.Message
		if len(args) > 0 {
			message = values.Message(strings.Join(args, " "))
		} else {
			var err error
			message, err = scan()
			if err != nil {
				return fmt.Errorf("scan message: %w", err)
			}
		}

		if err := mes.Send(ctx, message); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}

	return &Root{
		Command: rootCmd.Command,
	}
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

func scan() (values.Message, error) {
	sc := bufio.NewScanner(os.Stdin)
	sb := &strings.Builder{}
	for sc.Scan() {
		text := sc.Text()
		sb.WriteString(text + "\n")
	}

	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	return values.Message(strings.TrimSpace(sb.String())), nil
}
