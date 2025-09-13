/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cmp"
	"fmt"
	"os"
	"strings"

	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type InitBare struct {
	*cobra.Command
}

func NewInitBare(rootCmd *Root) *InitBare {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the configuration file",
		Long:  `initコマンドは、設定ファイルの値を対話形式で設定することができます。`,
	}

	rootCmd.AddCommand(initCmd)

	return &InitBare{
		Command: initCmd,
	}
}

type Init struct {
	*cobra.Command
}

// TODO: *viper.Viperじゃなくて設定ファイルに書き込む用の何かinterfaceを渡す
func NewInit(initBare *InitBare, initConfig config.Init, cw config.FileWriter) *Init {
	initBare.RunE = func(cmd *cobra.Command, args []string) error {
		force, err := initConfig.GetForce()
		if err != nil {
			return fmt.Errorf("failed to get force flag: %w", err)
		}

		filePath, err := cw.GetUsedFilePath()
		if err != nil {
			return fmt.Errorf("get config file path: %w", err)
		}
		if force {
			fmt.Printf("Overwriting the existing configuration file: %s\n", cmp.Or(filePath, "(config file not found)"))
		} else {
			if filePath != "" {
				fmt.Println("Configuration file already exists.")
				fmt.Println(filePath)
				return nil
			}
		}

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("terminal make raw: %w", err)
		}
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
			{"Enter the webhook host", "https://q.trap.jp", "webhook_host", false},
			{"Enter the webhook ID", "", "webhook_id", false},
			{"Enter the webhook secret", "", "webhook_secret", true},
		}

		t := term.NewTerminal(os.Stdin, "")
		values := make(map[string]any, len(prompts))
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
				return fmt.Errorf("terminal read line: %w", err)
			}

			input = strings.TrimSpace(input)
			if input == "" {
				if p.defaultValue == "" {
					cobra.CheckErr(fmt.Sprintf("%s is required", p.configKey))
				}
				input = p.defaultValue
			}

			values[p.configKey] = input
		}

		if err := cw.Write(force, config.ConfigValues(values)); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		return nil
	}

	return &Init{
		Command: initBare.Command,
	}
}
