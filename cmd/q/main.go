/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/ikura-hamu/q-cli/internal/cmd"
	"github.com/ikura-hamu/q-cli/internal/config/file"
	"github.com/ikura-hamu/q-cli/internal/config/flag"
	"github.com/ikura-hamu/q-cli/internal/message/impl"
	secretImpl "github.com/ikura-hamu/q-cli/internal/secret/impl"
)

func main() {
	rootBareCmd := cmd.NewRootBare()
	confFile := flag.NewFile(rootBareCmd.PersistentFlags())
	confRoot := flag.NewRoot(rootBareCmd.Flags())
	v, err := file.NewViper(confFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	confWebhook := file.NewWebhookFactory(v)
	clientFactory := webhook.NewWebhookClientFactory(confWebhook)
	if err != nil {
		fmt.Println("create client:", err)
		os.Exit(1)
	}
	mes := impl.NewMessage()
	sec := secretImpl.NewSecretDetector()

	rootCmd := cmd.NewRoot(rootBareCmd, confFile, confRoot, clientFactory, mes, sec)

	initBareCmd := cmd.NewInitBare(rootCmd)
	confInit := flag.NewInit(initBareCmd)
	configFileWriter := file.NewWriter(v)
	_ = cmd.NewInit(initBareCmd, confInit, configFileWriter)

	confBareCmd := cmd.NewConfigBare()
	configFileReader := file.NewReader(v)
	_ = cmd.NewConfig(rootCmd, confBareCmd, configFileReader)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
