package command

import (
	"context"
	"os"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/spf13/pflag"
)

type Base struct {
	Ui  cli.Ui
	cmd NamedCommand
}

func NewBase(ui cli.Ui, cmd NamedCommand) Base {
	return Base{
		Ui:  ui,
		cmd: cmd,
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
		return 1
	}

	return 0
}
