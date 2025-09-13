package file

import (
	"fmt"

	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/viper"
)

type Writer struct {
	v *viper.Viper
}

var _ config.FileWriter = (*Writer)(nil)

func NewWriter(v *viper.Viper) *Writer {
	return &Writer{
		v: v,
	}
}

func (w *Writer) GetUsedFilePath() (string, error) {
	return w.v.ConfigFileUsed(), nil
}

func (w *Writer) Write(force bool, val config.ConfigValues) error {
	for k, v := range val {
		w.v.Set(k, v)
	}

	if force {
		err := w.v.WriteConfig()
		if err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	} else {
		err := w.v.SafeWriteConfig()
		if err != nil {
			return fmt.Errorf("safe write config: %w", err)
		}
	}

	return nil
}
