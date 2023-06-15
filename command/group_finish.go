package command

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/codes"
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

func createSpan(tp trace.Tracer, traceParent string, props map[string]string) (trace.Span, error) {
	nano := props["start"]
	i, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, i)

	// cli carrier traceParent
	ctx := tracing.WithTraceParent(context.Background(), traceParent)
	_, span := tp.Start(ctx, props["name"], trace.WithTimestamp(start))

	return span, nil
}

func finishSpan(span trace.Span, finish int64) {
	span.End(trace.WithTimestamp(time.Unix(0, finish)))
}

func applyProps(span trace.Span, props map[string]string) {

	delete(props, "name")
	delete(props, "start")

	attrs := tracing.AttributesFromMap(props)

	span.SetAttributes(attrs...)
}

func applyStatus(span trace.Span, flag *pflag.Flag) {

	if flag != nil && flag.Changed {

		value := flag.Value.String()
		if value == "unset" {
			value = ""
		}

		span.SetStatus(codes.Error, value)
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
