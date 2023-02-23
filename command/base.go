package command

import (
	"context"
	"os"
	"os/exec"
	"path"
	"time"
	"trace/tracing"

	"github.com/go-logfmt/logfmt"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const TraceParentEnvVar = "TRACEPARENT"

type Base struct {
	Ui  cli.Ui
	cmd NamedCommand

	now              func() int64
	testSpanExporter tracesdk.SpanExporter
}

func NewBase(ui cli.Ui, cmd NamedCommand) Base {
	return Base{
		Ui:  ui,
		cmd: cmd,

		now: func() int64 { return time.Now().UnixNano() },
	}
}

type NamedCommand interface {
	Name() string
	Synopsis() string

	Flags() *pflag.FlagSet
	EnvironmentVariables() map[string]string

	RunContext(ctx context.Context, args []string) error
}

func (b *Base) AutocompleteFlags() complete.Flags {
	return nil
}

func (b *Base) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (b *Base) Help() string {
	return b.cmd.Synopsis() + "\n\n" + b.allFlags().FlagUsages()
}

func (b *Base) allFlags() *pflag.FlagSet {

	flags := b.cmd.Flags()

	return flags
}

func (b *Base) allEnvironmentVariables() map[string]string {

	vars := b.cmd.EnvironmentVariables()

	return vars
}

func (b *Base) applyEnvironmentFallback(flags *pflag.FlagSet) {
	envVars := b.allEnvironmentVariables()

	flags.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			return
		}

		envVar, found := envVars[f.Name]
		v := ""
		if found {
			v = os.Getenv(envVar)
		}

		isDifferent := v != f.DefValue

		if found && v != "" && isDifferent {
			f.Value.Set(v)
		}
	})
}

func (b *Base) Run(args []string) int {
	f := b.allFlags()

	if err := f.Parse(args); err != nil {
		b.Ui.Error(err.Error())
		return 1
	}

	b.applyEnvironmentFallback(f)

	if err := b.cmd.RunContext(context.Background(), f.Args()); err != nil {
		b.Ui.Error(err.Error())

		// handle an exitError specifically so we can pass its exit code on
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}

		return 1
	}

	return 0
}

func (b *Base) createTracer(ctx context.Context, g tracesdk.IDGenerator) (trace.Tracer, error) {

	exporter := b.testSpanExporter

	if exporter == nil {
		var err error
		exporter, err = tracing.CreateExporter(ctx, &tracing.ExporterConfig{
			Endpoint: "localhost:4317",
		})
		if err != nil {
			return nil, err
		}
	}

	tp, err := tracing.CreateTraceProvider(ctx, g, exporter)
	if err != nil {
		return nil, err
	}

	return tp.Tracer("trace-cli"), nil
}

func (b *Base) writeState(traceParent string, data map[string]any) error {

	dir := path.Join(os.TempDir(), "trace", "state")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(dir, traceParent), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := logfmt.NewEncoder(f)
	for key, val := range data {
		if err := encoder.EncodeKeyval(key, val); err != nil {
			return err
		}
	}

	if err := encoder.EndRecord(); err != nil {
		return err
	}

	return nil
}

func (b *Base) readState(traceParent string) (map[string]string, error) {

	filePath := path.Join(os.TempDir(), "trace", "state", traceParent)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := logfmt.NewDecoder(f)
	data := map[string]string{}

	for decoder.ScanRecord() {
		for decoder.ScanKeyval() {
			k := string(decoder.Key())
			v := string(decoder.Value())

			data[k] = v
		}
	}

	return data, nil
}
