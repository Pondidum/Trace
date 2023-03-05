package command

import (
	"context"
	"fmt"
	"strconv"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewStartCommand(ui cli.Ui) (*StartCommand, error) {
	cmd := &StartCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type StartCommand struct {
	Base

	attrPairs []string
}

func (c *StartCommand) Name() string {
	return "start"
}

func (c *StartCommand) Synopsis() string {
	return "Start a Trace"
}

func (c *StartCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringSliceVar(&c.attrPairs, "attr", []string{}, "")

	return flags
}

func (c *StartCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *StartCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("this command takes 1 argument: span_name")
	}

	name := args[0]
	traceParent := tracing.NewTraceParent()

	data, err := mapFromKeyValues(c.attrPairs)
	if err != nil {
		return err
	}

	data["name"] = name
	data["start"] = strconv.FormatInt(c.now(), 10)

	if err := c.writeState(traceParent, data); err != nil {
		return err
	}

	c.Ui.Output(traceParent)
	return nil
}
