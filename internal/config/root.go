package config

import "github.com/guregu/null/v6"

type Root interface {
	GetVersion() (bool, error)
	GetCodeBlock() (bool, error)
	GetCodeBlockLang() (null.String, error)
	GetChannelName() (null.String, error)
	GetPrintBeforeSend() (bool, error)
}
