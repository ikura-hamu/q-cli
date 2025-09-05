package flag

import (
	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/spf13/pflag"
)

type Root struct {
	version         bool
	codeBlock       bool
	codeBlockLang   string
	channelName     string
	printBeforeSend bool
}

var _ config.Root = (*Root)(nil)

func NewRoot(flagSet *pflag.FlagSet) *Root {
	r := &Root{}
	flagSet.BoolVarP(&r.version, "version", "v", false, "Print version information and exit.")
	flagSet.BoolVarP(&r.codeBlock, "code-block", "c", false, "Wrap the message in a code block.")
	flagSet.StringVarP(&r.codeBlockLang, "lang", "l", "", "Specify the language for the code block. Used only when --code-block is set.")
	flagSet.StringVarP(&r.channelName, "channel", "C", "", "Specify the channel name to send the message to. If not specified, the default channel will be used.")
	flagSet.BoolVarP(&r.printBeforeSend, "print-before-send", "p", false, "Print the message to be sent before sending it.")
	return r
}

func (r *Root) GetVersion() (bool, error) {
	return r.version, nil
}

func (r *Root) GetCodeBlock() (bool, error) {
	return r.codeBlock, nil
}

func (r *Root) GetCodeBlockLang() (null.String, error) {
	return null.NewString(r.codeBlockLang, r.codeBlockLang != ""), nil
}

func (r *Root) GetChannelName() (null.String, error) {
	return null.NewString(r.channelName, r.channelName != ""), nil
}

func (r *Root) GetPrintBeforeSend() (bool, error) {
	return r.printBeforeSend, nil
}
