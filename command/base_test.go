package command

import (
	"context"
	"testing"
	"time"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestBaseStruct(t *testing.T) {

	base := NewBase(cli.NewMockUi(), &FinishCommand{})

	start := base.now()

	time.Sleep(1 * time.Second)

	finish := base.now()

	assert.Greater(t, finish, start)

}

type MockCommand struct {
	Base
}

func (c *MockCommand) Name() string {
	return "mock"
}

func (c *MockCommand) Synopsis() string {
	return "mocks"
}

func (c *MockCommand) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
}

func (c *MockCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *MockCommand) RunContext(ctx context.Context, args []string) error {
	return nil
}
