package cmd

import (
	"fmt"
	"os"
)

func returnWithError(message string, args ...any) {
	fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(1)
}
