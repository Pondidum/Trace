package command

import (
	"context"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewGenerateCommand(ui cli.Ui) (*GenerateCommand, error) {
	cmd := &GenerateCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type GenerateCommand struct {
	Base
}

func (c *GenerateCommand) Name() string {
	return "generate"
}

func (c *GenerateCommand) Synopsis() string {
	return "Generate a Trace ID"
}

func (c *GenerateCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	return flags
}

func (c *GenerateCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *GenerateCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("this command takes 1 argument: span_name")
	}

	name := args[0]
	traceParent := NewTraceParent(ctx)

	err := c.writeState(traceParent, map[string]any{
		"name":  name,
		"start": c.now(),
	})

	if err != nil {
		return err
	}

	c.Ui.Output(traceParent)
	return nil
}
