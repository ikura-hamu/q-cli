package file

import (
	"fmt"
	"os"

	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/viper"
)

func NewViper(conf config.File) (*viper.Viper, error) {
	cfgFile, err := conf.GetFilePath()
	if err != nil {
		return nil, fmt.Errorf("get config file path: %w", err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(".q-cli")
	v.AutomaticEnv()
	v.SetEnvPrefix("Q")

	if cfgFile.Valid {
		v.SetConfigFile(cfgFile.ValueOrZero())
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home directory: %w", err)
		}

		// Search config in home directory with name ".q-cli" (without extension).
		v.AddConfigPath(home)
	}

	return v, nil
}
