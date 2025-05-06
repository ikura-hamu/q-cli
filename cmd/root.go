/*
Copyright © 2024 ikura-hamu 104292023+ikura-hamu@users.noreply.github.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/ikura-hamu/q-cli/internal/message"
	"github.com/ikura-hamu/q-cli/internal/message/impl"
	"github.com/ikura-hamu/q-cli/internal/secret"
	secretImpl "github.com/ikura-hamu/q-cli/internal/secret/impl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	cfgFile string

	cl             client.Client
	mes            message.Message
	secretDetector secret.SecretDetector
)

var (
	ErrEmptyConfiguration = errors.New("some webhook configuration field(s) is empty")
	ErrChannelNotFound    = errors.New("channel name is not in the configuration")
)

func SetClient(c client.Client) {
	cl = c
}

var Command = rootCmd

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	TraverseChildren: true,
	Use:              "q [message]",
	Example:          "q print(\"Hello, world!\") -c -l py",
	Short:            "traQ Webhook CLI",
	Long: `"q-cli" は、traQにWebhookを使ってメッセージを送信するためのCLIツールです。設定に基づいてWebhookを送信します。
設定は設定ファイルに記述するか、環境変数で指定することができます。
メッセージは標準入力からも受け取ることができます。
設定ファイルの場所は、何も指定しない場合、$HOME/.q-cli.yaml です。`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cl, err := webhook.NewWebhookClient(viper.GetString(configKeyWebhookID), viper.GetString(configKeyWebhookHost), viper.GetString(configKeyWebhookSecret))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		SetClient(cl)

		mes = impl.NewMessage()
		secretDetector = secretImpl.NewSecretDetector()

		return nil
	},

	// Uncomment the following line if your bare application
	// has an action associated with it
	RunE: func(cmd *cobra.Command, args []string) error {
		if printVersion {
			printVersionInfo()
			return nil
		}

		cmdCtx := cmd.Context()
		rootFlagsData, ok := cmdCtx.Value(rootFlagsCtxKey{}).(*rootFlags)
		if !ok {
			return errors.New("failed to get root options")
		}

		if cl == nil || mes == nil || secretDetector == nil {
			return errors.New("client, message or secret detector is nil")
		}

		channelsStr := viper.GetStringMapString(configKeyChannels)
		channels := make(map[string]uuid.UUID, len(channelsStr))
		for k, v := range channelsStr {
			id, err := uuid.Parse(v)
			if err != nil {
				return fmt.Errorf("failed to parse channel ID: %w", err)
			}
			channels[k] = id
		}

		conf := webhookConfig{
			host:     viper.GetString(configKeyWebhookHost),
			id:       viper.GetString(configKeyWebhookID),
			secret:   viper.GetString(configKeyWebhookSecret),
			channels: channels,
		}

		if conf.host == "" || conf.id == "" || conf.secret == "" {
			return ErrEmptyConfiguration
		}

		messageStr, err := mes.BuildMessage(args, message.Option{
			CodeBlock:     rootFlagsData.codeBlock,
			CodeBlockLang: rootFlagsData.codeBlockLang,
		})
		if err != nil {
			return fmt.Errorf("failed to build message: %w", err)
		}

		err = secretDetector.Detect(cmdCtx, messageStr)
		if detectMes, ok := secret.SecretDetected(err); ok {
			fmt.Println(detectMes)
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to detect secret: %w", err)
		}

		channelID := uuid.Nil
		if rootFlagsData.channelName != "" {
			var ok bool
			channelID, ok = conf.channels[rootFlagsData.channelName]
			if !ok {
				return fmt.Errorf("channel '%s' is not found: %w", rootFlagsData.channelName, ErrChannelNotFound)
			}
		}

		if rootFlagsData.printBeforeSend {
			ok, err := checkMessage(messageStr)
			if err != nil {
				return fmt.Errorf("failed to check message: %w", err)
			}
			if !ok {
				fmt.Println("Send canceled.")
				return nil
			}
		}

		if cl != nil {
			err := cl.SendMessage(messageStr, channelID)
			if errors.Is(err, client.ErrEmptyMessage) {
				return errors.New("empty message is not allowed")
			}
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		} else {
			panic("client is nil")
		}

		return nil
	},
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
	t.SetPrompt("Send? y/n(any): ")
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

type rootFlags struct {
	codeBlock       bool
	codeBlockLang   string
	channelName     string
	printBeforeSend bool
}

var (
	printVersion bool
	version      string
)

const (
	configKeyWebhookHost   = "webhook_host"
	configKeyWebhookID     = "webhook_id"
	configKeyWebhookSecret = "webhook_secret"
	configKeyChannels      = "channels"
)

type webhookConfig struct {
	host     string
	id       string
	secret   string
	channels map[string]uuid.UUID
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type rootFlagsCtxKey struct{}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "設定ファイルの場所を指定します。 (デフォルトは $HOME/.q-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version information and quit")

	var rootFlagsData rootFlags
	rootCmd.Flags().BoolVarP(&rootFlagsData.codeBlock, "code", "c", false, "コードブロック付きでメッセージを送信します。")
	rootCmd.Flags().StringVarP(&rootFlagsData.codeBlockLang, "lang", "l", "", "コードブロックの言語を指定します。")
	rootCmd.Flags().StringVarP(&rootFlagsData.channelName, "channel", "C", "", `チャンネル名を指定して、デフォルト以外のチャンネルにメッセージを送信します。
チャンネル名は設定ファイルの channels に記述されたキーを指定します。`)
	rootCmd.Flags().BoolVarP(&rootFlagsData.printBeforeSend, "print", "p", false, "メッセージを送信する前に表示し、確認を求めます。")

	ctx := context.WithValue(context.Background(), rootFlagsCtxKey{}, &rootFlagsData)
	rootCmd.SetContext(ctx)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".q-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".q-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetEnvPrefix("q")

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
	_ = viper.ReadInConfig()
}

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
