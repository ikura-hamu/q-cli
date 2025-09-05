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
	rootCmd := cmd.New()
	confFile := flag.NewFile(rootCmd.PersistentFlags())
	confRoot := flag.NewRoot(rootCmd.Flags())
	v, err := file.NewViper(confFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	confWebhook := file.NewWebhook(v)
	cl, err := webhook.NewClientFromConfig(confWebhook)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mes := impl.NewMessage()
	sec := secretImpl.NewSecretDetector()

	rootCmd = cmd.NewRoot(rootCmd, confFile, confRoot, confWebhook, cl, mes, sec)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
