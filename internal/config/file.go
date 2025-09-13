package config

import "github.com/guregu/null/v6"

type File interface {
	GetFilePath() (null.String, error)
}

type ConfigValues map[string]any

type FileNameGetter interface {
	GetUsedFilePath() (string, error)
}

type FileWriter interface {
	FileNameGetter
	Write(force bool, val ConfigValues) error
}

type FileReader interface {
	FileNameGetter
	Read() (ConfigValues, error)
}
