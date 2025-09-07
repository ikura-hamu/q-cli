package flag

import (
	"github.com/ikura-hamu/q-cli/internal/cmd"
	"github.com/ikura-hamu/q-cli/internal/config"
)

type Init struct {
	force bool
}

var _ config.Init = (*Init)(nil)

func NewInit(c *cmd.InitBare) *Init {
	init := &Init{}
	c.Flags().BoolVarP(&init.force, "force", "f", false, "既存の設定ファイルを上書きします")

	return init
}

func (i *Init) GetForce() (bool, error) {
	return i.force, nil
}
