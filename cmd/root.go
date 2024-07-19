/*
Copyright © 2024 ikura-hamu 104292023+ikura-hamu@users.noreply.github.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var cl client.Client

func SetClient(c client.Client) {
	cl = c
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "q-cli",
	Short: "traQ Webhook CLI",
	Long: `"q-cli" is a CLI tool for sending messages to traQ via webhook.
It reads the configuration file and sends the message to the specified webhook.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			printVersionInfo()
			return
		}

		conf := webhookConfig{
			host:   viper.GetString(configKeyWebhookHost),
			id:     viper.GetString(configKeyWebhookID),
			secret: viper.GetString(configKeyWebhookSecret),
		}

		if conf.host == "" || conf.id == "" || conf.secret == "" {
			returnWithError("some webhook configuration field(s) is empty\n")
			return
		}
		var message string

		if len(args) > 0 {
			message = strings.Join(args, " ")
		} else {
			sc := bufio.NewScanner(os.Stdin)
			sb := &strings.Builder{}
			for sc.Scan() {
				text := sc.Text()
				sb.WriteString(text + "\n")
			}
			message = sb.String()
		}

		if withCodeBlock {
			message = fmt.Sprintf("```%s\n%s\n```", codeBlockLang, message)
		}

		if cl != nil {

			err := cl.SendMessage(message)
			// err := makeWebhookRequest(conf, message)
			if err != nil {
				returnWithError("failed to send message: %v\n", err)
			}
		} else {
			panic("client is nil")
		}

		return
	},
}

var (
	printVersion  bool
	withCodeBlock bool
	codeBlockLang string
	version       string
)

const (
	configKeyWebhookHost   = "webhook_host"
	configKeyWebhookID     = "webhook_id"
	configKeyWebhookSecret = "webhook_secret"
)

type webhookConfig struct {
	host   string
	id     string
	secret string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.q-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version information and quit")

	rootCmd.Flags().BoolVarP(&withCodeBlock, "code", "c", false, "Send message with code block")
	rootCmd.Flags().StringVarP(&codeBlockLang, "lang", "l", "text", "Code block language")
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
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	cl, err := webhook.NewWebhookClient(viper.GetString(configKeyWebhookID), viper.GetString(configKeyWebhookHost), viper.GetString(configKeyWebhookSecret))
	if err != nil {
		returnWithError("failed to create webhook client: %v\n", err)
	}
	SetClient(cl)
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
