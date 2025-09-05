package config

import "github.com/guregu/null/v6"

type File interface {
	GetFilePath() (null.String, error)
}
