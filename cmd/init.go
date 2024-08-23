/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cmp"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
	Long:  `This command initializes the configuration file of q-cli interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		if initForce {
			fmt.Printf("Overwriting the existing configuration file: %s\n", cmp.Or(viper.ConfigFileUsed(), "(config file not found)"))
		} else {
			if viper.ConfigFileUsed() != "" {
				fmt.Println("Configuration file already exists.")
				fmt.Println(viper.ConfigFileUsed())
				return
			}
		}

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		cobra.CheckErr(err)
		defer func() {
			err := term.Restore(int(os.Stdin.Fd()), oldState)
			cobra.CheckErr(err)
		}()

		prompts := []struct {
			prompt       string
			defaultValue string // if not set, the input value is required
			configKey    string
			isPassword   bool
		}{
			{"Enter the webhook host", "https://q.trap.jp", configKeyWebhookHost, false},
			{"Enter the webhook ID", "", configKeyWebhookID, false},
			{"Enter the webhook secret", "", configKeyWebhookSecret, true},
		}

		t := term.NewTerminal(os.Stdin, "")
		for _, p := range prompts {
			prompt := p.prompt
			if p.defaultValue != "" {
				prompt += fmt.Sprintf(" (default: %s)", p.defaultValue)
			}
			prompt = fmt.Sprintf("%s: ", prompt)

			var input string
			var err error
			if p.isPassword {
				input, err = t.ReadPassword(prompt)
			} else {
				t.SetPrompt(prompt)
				input, err = t.ReadLine()
			}
			if err != nil {
				cobra.CheckErr(err)
			}

			input = strings.TrimSpace(input)
			if input == "" {
				if p.defaultValue == "" {
					cobra.CheckErr(fmt.Sprintf("%s is required", p.configKey))
				}
				input = p.defaultValue
			}

			viper.Set(p.configKey, input)
		}

		if initForce {
			err = viper.WriteConfig()
		} else {
			err = viper.SafeWriteConfig()
		}
		cobra.CheckErr(err)
	},
}

var (
	initForce bool
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "If provided, overwrite the existing configuration file")
}
