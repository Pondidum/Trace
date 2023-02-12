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
	"go.opentelemetry.io/otel/trace"
)

func NewSpanFinishCommand(ui cli.Ui) (*SpanFinishCommand, error) {
	cmd := &SpanFinishCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type SpanFinishCommand struct {
	Base
}

func (c *SpanFinishCommand) Name() string {
	return "span finish"
}

func (c *SpanFinishCommand) Synopsis() string {
	return "Finish a span"
}

func (c *SpanFinishCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	return flags
}

func (c *SpanFinishCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *SpanFinishCommand) RunContext(ctx context.Context, args []string) error {

	spanid := os.Getenv(TraceParentEnvVar)
	if len(args) > 0 {
		spanid = args[0]
	}

	if spanid == "" {
		return fmt.Errorf("this command takes 1 argument: spanid")
	}

	state, err := c.readState(spanid)
	if err != nil {
		return err
	}

	traceId, _, err := tracing.ParseTraceParent(spanid)
	if err != nil {
		return err
	}

	parentSpan, err := trace.SpanIDFromHex(state["parent"])
	if err != nil {
		return err
	}

	parentSpanId := AsTraceParent(traceId, parentSpan)

	ids, err := tracing.ContinueExisting(spanid)
	if err != nil {
		return err
	}

	tracer, err := c.createTracer(ctx, ids)
	if err != nil {
		return err
	}

	if err := createSpan(tracer, parentSpanId, c.now(), state); err != nil {
		return err
	}

	return nil
}

func createSpan(tp trace.Tracer, traceParent string, finish int64, props map[string]string) error {

	nano := props["start"]
	i, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return err
	}
	start := time.Unix(0, i)

	// cli carrier traceParent
	ctx := tracing.WithTraceParent(context.Background(), traceParent)
	_, span := tp.Start(ctx, props["name"], trace.WithTimestamp(start))

	delete(props, "name")
	delete(props, "start")
	attrs := tracing.FromMap(props)
	span.SetAttributes(attrs...)

	span.End(trace.WithTimestamp(time.Unix(0, finish)))

	return nil
}
