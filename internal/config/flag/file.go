package flag

import (
	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/pflag"
)

type File struct {
	filePath string
}

var _ config.File = (*File)(nil)

func NewFile(flagSet *pflag.FlagSet) *File {
	f := &File{}
	flagSet.StringVar(&f.filePath, "config", "", "Config file path. (default is $HOME/.q-cli.yaml)")
	return f
}

func (f *File) GetFilePath() (null.String, error) {
	if f.filePath == "" {
		return null.String{}, nil
	}
	return null.StringFrom(f.filePath), nil
}
