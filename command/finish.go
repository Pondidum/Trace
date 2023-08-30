package command

import (
	"context"
	"fmt"
	"os"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewFinishCommand(ui cli.Ui) (*FinishCommand, error) {
	cmd := &FinishCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type FinishCommand struct {
	Base

	attrPairs []string
}

func (c *FinishCommand) Name() string {
	return "Finish"
}

func (c *FinishCommand) Synopsis() string {
	return "Finish a Trace ID"
}

func (c *FinishCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringSliceVar(&c.attrPairs, "attr", []string{}, "")

	flags.String("error", "", "")
	flags.Lookup("error").NoOptDefVal = "unset"

	return flags
}

func (c *FinishCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *FinishCommand) RunContext(ctx context.Context, args []string) error {

	traceParent := os.Getenv(TraceParentEnvVar)
	if len(args) > 0 {
		traceParent = args[0]
	}

	if traceParent == "" {
		return fmt.Errorf("this command takes 1 argument: traceparent")
	}

	ids, err := tracing.ContinueExisting(traceParent)
	if err != nil {
		return err
	}

	state, err := c.readState(traceParent)
	if err != nil {
		return err
	}

	attrs, err := mapFromKeyValues(c.attrPairs)
	if err != nil {
		return err
	}

	for k, v := range attrs {
		state[k] = v
	}

	tracer, err := c.createTracer(ctx, ids)
	if err != nil {
		return err
	}

	span, err := createRootSpan(tracer, traceParent, state)
	if err != nil {
		return err
	}

	applyProps(span, state)
	applyStatus(span, c.allFlags().Lookup("error"))
	finishSpan(span, c.now())

	return nil
}
