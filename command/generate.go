package command

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"
)

func NewGenerateCommand(ui cli.Ui) (*GenerateCommand, error) {
	cmd := &GenerateCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type GenerateCommand struct {
	Base
}

func (c *GenerateCommand) Name() string {
	return "generate"
}

func (c *GenerateCommand) Synopsis() string {
	return "Generate a Trace ID"
}

func (c *GenerateCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	return flags
}

func (c *GenerateCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *GenerateCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("this command takes 1 argument: span_name")
	}

	name := args[0]

	generator := newIDGenerator()
	trace, span := generator.NewIDs(ctx)

	traceParent := fmt.Sprintf("00-%s-%s-01", trace, span)

	err := c.writeState(traceParent, map[string]any{
		"name":  name,
		"start": c.now(),
	})

	if err != nil {
		return err
	}

	c.Ui.Output(traceParent)
	return nil
}

type randomIDGenerator struct {
	randSource *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	tid := trace.TraceID{}
	_, _ = gen.randSource.Read(tid[:])
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return tid, sid
}

func newIDGenerator() *randomIDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))

	return gen
}
