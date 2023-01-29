package command

import (
	"github.com/mitchellh/cli"
)

func Commands(ui cli.Ui) map[string]cli.CommandFactory {

	return map[string]cli.CommandFactory{
		"generate": func() (cli.Command, error) {
			return NewGenerateCommand(ui)
		},
		"finish": func() (cli.Command, error) {
			return NewFinishCommand(ui)
		},
	}
}
