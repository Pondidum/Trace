package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewTaskCommand(ui cli.Ui) (*TaskCommand, error) {
	cmd := &TaskCommand{}
	cmd.Base = NewBase(ui, cmd)

	return cmd, nil
}

type TaskCommand struct {
	Base

	taskName string
}

func (c *TaskCommand) Name() string {
	return "task"
}

func (c *TaskCommand) Synopsis() string {
	return "Runs a command, storing exit code etc."
}

func (c *TaskCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringVar(&c.taskName, "name", "", "name the task running")
	return flags
}

func (c *TaskCommand) EnvironmentVariables() map[string]string {
	return map[string]string{}
}

func (c *TaskCommand) RunContext(ctx context.Context, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("this command takes at least 1 argument")
	}

	traceParent := args[0]
	commandAndArgs := args[1:]

	if _, _, err := tracing.ParseTraceParent(traceParent); err != nil {

		fromEnv := os.Getenv(TraceParentEnvVar)

		if fromEnv == "" {
			return fmt.Errorf("no trace parent from env and first arg wasnt trace")

		} else if _, _, err := tracing.ParseTraceParent(fromEnv); err != nil {
			return fmt.Errorf(" trace parent from env was invalid")
		}

		traceParent = fromEnv
		commandAndArgs = args
	}

	if len(commandAndArgs) == 0 {
		return fmt.Errorf("you must specify a command")
	}

	var quoted []string
	for _, s := range commandAndArgs {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", strings.Replace(s, "\"", "\\\"", -1)))
	}

	script := strings.Join(quoted, " ")

	startTime := c.now()

	taskError := runTask(script)

	finishTime := c.now()

	tracer, err := c.createTracer(ctx, nil)
	if err != nil {
		return err
	}

	taskName := c.buildTaskName(commandAndArgs)

	spanContext := tracing.WithTraceParent(context.Background(), traceParent)
	_, span := tracer.Start(spanContext, taskName, trace.WithTimestamp(time.Unix(0, startTime)))

	span.SetAttributes(
		attribute.String("command.executable", commandAndArgs[0]),
		attribute.StringSlice("command.arguments", commandAndArgs[1:]),
		attribute.String("shell_script", script),
		attribute.Int("command.exit_code", 0), // assumed for now, replaced later
	)

	if taskError == nil {
		span.SetStatus(codes.Ok, "")
	} else {
		span.SetStatus(codes.Error, taskError.Error())

		if exitErr, ok := err.(*exec.ExitError); ok {
			span.SetAttributes(attribute.Int("command.exit_code", exitErr.ExitCode()))
		}
	}

	span.End(trace.WithTimestamp(time.Unix(0, finishTime)))

	return taskError
}

func (c *TaskCommand) buildTaskName(command []string) string {

	if strings.TrimSpace(c.taskName) != "" {
		return c.taskName
	}

	return strings.Join(command, " ")
}

func runTask(script string) error {

	shell := "/bin/sh"
	if val := os.Getenv("SHELL"); val != "" {
		shell = val
	}

	cmd := exec.Command(shell, "-c", script)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
