/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/webhook"
	"github.com/ikura-hamu/q-cli/internal/cmd"
	"github.com/ikura-hamu/q-cli/internal/config/file"
	"github.com/ikura-hamu/q-cli/internal/config/flag"
	"github.com/ikura-hamu/q-cli/internal/interaction/terminal"
	secretImpl "github.com/ikura-hamu/q-cli/internal/secret/impl"
	serviceImpl "github.com/ikura-hamu/q-cli/internal/service/impl"
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
	sec := secretImpl.NewSecretDetector()
	termSess := terminal.NewSession()

	mesService := serviceImpl.NewMessage(clientFactory, confRoot, sec, termSess)

	rootCmd := cmd.NewRoot[client.Client](rootBareCmd, confRoot, mesService)

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
