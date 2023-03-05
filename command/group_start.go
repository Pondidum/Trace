package command

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewGroupStartCommand(ui cli.Ui) (*GroupStartCommand, error) {
	cmd := &GroupStartCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type GroupStartCommand struct {
	Base

	attrPairs []string
}

func (c *GroupStartCommand) Name() string {
	return "group start"
}

func (c *GroupStartCommand) Synopsis() string {
	return "start a group"
}

func (c *GroupStartCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	flags.StringSliceVar(&c.attrPairs, "attr", []string{}, "")

	return flags
}

func (c *GroupStartCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *GroupStartCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("this command takes at least 1 argument: group_name")
	}

	name := args[0]
	traceParent := os.Getenv(TraceParentEnvVar)
	if len(args) > 1 {
		traceParent = args[1]
	}

	if traceParent == "" {
		return fmt.Errorf("this command requires a trace_parent, either from the command line or environment")
	}

	tid, parentSid, err := tracing.ParseTraceParent(traceParent)
	if err != nil {
		return err
	}

	sid := tracing.NewSpanID()
	newTraceParent := tracing.AsTraceParent(tid, sid)

	data, err := mapFromKeyValues(c.attrPairs)
	if err != nil {
		return err
	}

	data["name"] = name
	data["start"] = strconv.FormatInt(c.now(), 10)
	data["parent"] = parentSid.String()

	if err := c.writeState(newTraceParent, data); err != nil {
		return err
	}

	c.Ui.Output(newTraceParent)

	return nil
}
