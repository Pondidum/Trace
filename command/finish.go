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

func NewFinishCommand(ui cli.Ui) (*FinishCommand, error) {
	cmd := &FinishCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type FinishCommand struct {
	Base
}

func (c *FinishCommand) Name() string {
	return "Finish"
}

func (c *FinishCommand) Synopsis() string {
	return "Finish a Trace ID"
}

func (c *FinishCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
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

	state, err := c.readState(traceParent)
	if err != nil {
		return err
	}

	ids, err := tracing.ContinueExisting(traceParent)
	if err != nil {
		return err
	}

	tracer, err := c.createTracer(ctx, ids)
	if err != nil {
		return err
	}

	if err := createRootSpan(tracer, traceParent, c.now(), state); err != nil {
		return err
	}

	return nil
}

func createRootSpan(tp trace.Tracer, traceParent string, finish int64, props map[string]string) error {

	nano := props["start"]
	i, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return err
	}
	start := time.Unix(0, i)

	_, span := tp.Start(context.Background(), props["name"], trace.WithTimestamp(start))

	delete(props, "name")
	delete(props, "start")
	attrs := tracing.FromMap(props)
	span.SetAttributes(attrs...)

	span.End(trace.WithTimestamp(time.Unix(0, finish)))

	return nil
}
