package command

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

const ISO8601 = "2006-01-02T15:04:05-0700"

func NewStartCommand(ui cli.Ui) (*StartCommand, error) {
	cmd := &StartCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type StartCommand struct {
	Base

	attrPairs []string
	startTime string
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
	flags.StringVar(&c.startTime, "when", "", "ISO 8601 or epoch formatted time representing when the span starts")

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

	startEpoch, err := c.startEpochNano()
	if err != nil {
		return err
	}

	data["name"] = name
	data["start"] = strconv.FormatInt(startEpoch, 10)

	if err := c.writeState(traceParent, data); err != nil {
		return err
	}

	c.Ui.Output(traceParent)
	return nil
}

func (c *StartCommand) startEpochNano() (int64, error) {

	if c.startTime == "" {
		return c.now(), nil
	}

	seconds, err := strconv.ParseInt(c.startTime, 10, 64)
	if err == nil {
		return time.Unix(seconds, 0).UnixNano(), nil
	}

	t, err := time.Parse(ISO8601, c.startTime)
	if err != nil {
		return 0, err
	}

	return t.UnixNano(), nil
}
