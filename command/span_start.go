package command

import (
	"context"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewSpanStartCommand(ui cli.Ui) (*SpanStartCommand, error) {
	cmd := &SpanStartCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type SpanStartCommand struct {
	Base
}

func (c *SpanStartCommand) Name() string {
	return "span start"
}

func (c *SpanStartCommand) Synopsis() string {
	return "start a span"
}

func (c *SpanStartCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	return flags
}

func (c *SpanStartCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *SpanStartCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("this command takes at least 1 argument: span_name")
	}

	name := args[0]
	traceParent := os.Getenv(TraceParentEnvVar)
	if len(args) > 1 {
		traceParent = args[1]
	}

	if traceParent == "" {
		return fmt.Errorf("this command requires a trace_parent, either from the command line or environment")
	}

	tid, parentSid, err := ParseTraceParent(traceParent)
	if err != nil {
		return err
	}

	sid := NewSpanID()
	newTraceParent := AsTraceParent(tid, sid)

	data := map[string]any{
		"name":   name,
		"start":  c.now(),
		"parent": parentSid,
	}

	if err := c.writeState(newTraceParent, data); err != nil {
		return err
	}

	c.Ui.Output(newTraceParent)

	return nil
}
