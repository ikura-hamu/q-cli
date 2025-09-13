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
	confWebhook := file.NewWebhook(v)
	clientFactory, err := webhook.NewClientFromConfig(confWebhook)
	if err != nil {
		fmt.Println("create client:", err)
		os.Exit(1)
	}
	mes := impl.NewMessage()
	sec := secretImpl.NewSecretDetector()

	rootCmd := cmd.NewRoot(rootBareCmd, confFile, confRoot, clientFactory, mes, sec, v)

	initBareCmd := cmd.NewInitBare(rootCmd)
	confInit := flag.NewInit(initBareCmd)
	_ = cmd.NewInit(initBareCmd, confInit, v)

	confBareCmd := cmd.NewConfigBare()
	_ = cmd.NewConfig(rootCmd, confBareCmd, confWebhook, v)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
