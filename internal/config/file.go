package config

import "github.com/guregu/null/v6"

type File interface {
	GetFilePath() (null.String, error)
}

type ConfigValues map[string]any

type FileWriter interface {
	GetUsedFilePath() (string, error)
	Write(force bool, val ConfigValues) error
}
