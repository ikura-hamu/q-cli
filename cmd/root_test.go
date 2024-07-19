package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestRoot(t *testing.T) {
	t.Run("rootCmd", func(t *testing.T) {
		t.Log("host", viper.GetString("webhook_host"))
		t.Log("version", printVersion)

		rootCmd.DebugFlags()
		// rootCmd.Flag("version").Value.Set("true")
		rootCmd.DebugFlags()
		viper.Set("webhook_host", "http://localhost:8080")

		t.Log("version", printVersion)
		t.Log(cl)

		rootCmd.Run(rootCmd, []string{"test"})

	})
}
