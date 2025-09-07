/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ConfigBare struct {
	*cobra.Command
}

func NewConfigBare() *ConfigBare {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "config prints the configuration file path and the current configuration of the CLI",
		Long:  `configコマンドは、設定ファイルのパスと現在のCLIの設定を表示します。webhook_secretなど、一部の設定はマスクされます。`,
	}
	return &ConfigBare{
		Command: configCmd,
	}
}

type Config struct {
	*cobra.Command
}

func NewConfig(root *Root, confBareCmd *ConfigBare, conf config.Webhook, v *viper.Viper) *Config {
	root.AddCommand(confBareCmd.Command)
	confBareCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		err := v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}
		return nil
	}

	confBareCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fileName := v.ConfigFileUsed()

		allConfig := v.AllSettings()

		allConfig["webhook_secret"] = "********"

		yamlConfig, err := yaml.Marshal(allConfig)
		cobra.CheckErr(err)

		fmt.Printf(`Config file used: %s

%s
`, fileName, string(yamlConfig))

		return nil
	}
	return &Config{
		Command: confBareCmd.Command,
	}
}
