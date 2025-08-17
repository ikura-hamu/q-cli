/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config prints the configuration file path and the current configuration of the CLI",
	Long:  `configコマンドは、設定ファイルのパスと現在のCLIの設定を表示します。webhook_secretなど、一部の設定はマスクされます。`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := viper.ConfigFileUsed()
		allConfig := viper.AllSettings()

		allConfig[configKeyWebhookSecret] = "********"

		yamlConfig, err := yaml.Marshal(allConfig)
		cobra.CheckErr(err)

		fmt.Printf(`Config file used: %s

%s
`, fileName, string(yamlConfig))

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
