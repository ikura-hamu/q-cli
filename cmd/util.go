package cmd

import (
	"fmt"
	"os"
)

func printError(message string, args ...any) {
	fmt.Fprintf(os.Stderr, "q-cli error: "+message, args...)
}
