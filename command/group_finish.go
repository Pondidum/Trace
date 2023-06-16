package command

import (
	"context"
	"fmt"
	"os"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"
)

func NewGroupFinishCommand(ui cli.Ui) (*GroupFinishCommand, error) {
	cmd := &GroupFinishCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type GroupFinishCommand struct {
	Base

	attrPairs []string
}

func (c *GroupFinishCommand) Name() string {
	return "group finish"
}

func (c *GroupFinishCommand) Synopsis() string {
	return "Finish a group"
}

func (c *GroupFinishCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringSliceVar(&c.attrPairs, "attr", []string{}, "")

	flags.String("error", "", "")
	flags.Lookup("error").NoOptDefVal = "unset"

	return flags
}

func (c *GroupFinishCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *GroupFinishCommand) RunContext(ctx context.Context, args []string) error {

	groupId := os.Getenv(TraceParentEnvVar)
	if len(args) > 0 {
		groupId = args[0]
	}

	if groupId == "" {
		return fmt.Errorf("this command takes 1 argument: groupid")
	}

	traceId, _, err := tracing.ParseTraceParent(groupId)
	if err != nil {
		return err
	}

	props, err := c.readState(groupId)
	if err != nil {
		return err
	}

	attrs, err := mapFromKeyValues(c.attrPairs)
	if err != nil {
		return err
	}

	for k, v := range attrs {
		props[k] = v
	}

	parentSpan, err := trace.SpanIDFromHex(props["parent"])
	if err != nil {
		return err
	}

	parentSpanId := tracing.AsTraceParent(traceId, parentSpan)

	ids, err := tracing.ContinueExisting(groupId)
	if err != nil {
		return err
	}

	tracer, err := c.createTracer(ctx, ids)
	if err != nil {
		return err
	}

	span, err := createSpan(tracer, parentSpanId, props)
	if err != nil {
		return err
	}

	applyProps(span, props)
	applyStatus(span, c.allFlags().Lookup("error"))
	finishSpan(span, c.now())

	return nil
}
