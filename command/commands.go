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

		"span start": func() (cli.Command, error) {
			return NewSpanStartCommand(ui)
		},
		"span finish": func() (cli.Command, error) {
			return NewSpanFinishCommand(ui)
		},
	}
}
