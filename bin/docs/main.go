package main

import (
	"log"

	"github.com/ikura-hamu/q-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	if err := doc.GenMarkdownTree(cmd.Command, "docs"); err != nil {
		log.Fatalln(err)
	}
}
