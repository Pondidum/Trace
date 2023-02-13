package command

import (
	"github.com/mitchellh/cli"
)

func Commands(ui cli.Ui) map[string]cli.CommandFactory {

	return map[string]cli.CommandFactory{
		"start": func() (cli.Command, error) {
			return NewStartCommand(ui)
		},
		"finish": func() (cli.Command, error) {
			return NewFinishCommand(ui)
		},

		"group": func() (cli.Command, error) {
			return &cli.MockCommand{
				SynopsisText: "Interact with groups",
				HelpText:     "Interact with groups",
				RunResult:    cli.RunResultHelp,
			}, nil
		},
		"group start": func() (cli.Command, error) {
			return NewGroupStartCommand(ui)
		},
		"group finish": func() (cli.Command, error) {
			return NewGroupFinishCommand(ui)
		},

		"version": func() (cli.Command, error) {
			return NewVersionCommand(ui)
		},
	}
}
