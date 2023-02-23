package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

func NewAttrCommand(ui cli.Ui) (*AttrCommand, error) {
	cmd := &AttrCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type AttrCommand struct {
	Base
}

func (c *AttrCommand) Name() string {
	return "attr"
}

func (c *AttrCommand) Synopsis() string {
	return "Add attributes to a trace or span"
}

func (c *AttrCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	return flags
}

func (c *AttrCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *AttrCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("this command takes at least 1 argument: groupid")
	}

	groupId := args[0]
	pairs := args[1:]

	if _, _, err := tracing.ParseTraceParent(groupId); err != nil {

		fromEnv := os.Getenv(TraceParentEnvVar)

		if fromEnv == "" {
			return fmt.Errorf("the first argument wasn't a valid traceParent, and the $%s envvar was not specified", TraceParentEnvVar)

		} else if _, _, err := tracing.ParseTraceParent(fromEnv); err != nil {
			return fmt.Errorf("the traceParent from $%s was invalid: %w", TraceParentEnvVar, err)
		}

		groupId = fromEnv
		pairs = args
	}

	data, err := mapFromKeyValues(pairs)
	if err != nil {
		return err
	}

	if err := c.writeState(groupId, data); err != nil {
		return err
	}

	return nil
}

func mapFromKeyValues(pairs []string) (map[string]any, error) {
	m := make(map[string]any, len(pairs))

	for _, pair := range pairs {
		split := strings.Split(pair, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("'%s' was not a valid key=value pair", pair)
		}

		k := split[0]
		v := split[1]

		if k == "" || v == "" {
			return nil, fmt.Errorf("'%s' was not a valid key=value pair", pair)
		}

		m[k] = v
	}

	return m, nil
}
