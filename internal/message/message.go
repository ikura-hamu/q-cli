package message

type Option struct {
	CodeBlock     bool
	CodeBlockLang string
}

type Message interface {
	BuildMessage(args []string, option Option) string
	ContainsSecret(message string) (bool, string, error)
}
