package command

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestGeneration(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	cmd, _ := NewGenerateCommand(ui)

	assert.Equal(t, 0, cmd.Run([]string{}))
	assert.Regexp(t, `^[[:xdigit:]]{2}-[[:xdigit:]]{32}-[[:xdigit:]]{16}-[[:xdigit:]]{2}\n`, ui.OutputWriter.String())
}
