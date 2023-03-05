package command

import (
	"os"
	"testing"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestAttributes(t *testing.T) {

	cases := []struct {
		description   string
		traceParent   string
		args          []string
		expectedExit  int
		expectedError string
		expectedAttrs map[string]string
	}{
		{
			description:   "no args and no environment",
			expectedExit:  1,
			expectedError: "this command takes at least 1 argument: groupid",
		},
		{
			description:   "only keyvaluepair, no trace from env",
			args:          []string{"one=true"},
			expectedExit:  1,
			expectedError: "the first argument wasn't a valid traceParent, and the $TRACEPARENT envvar was not specified",
		},
		{
			description: "keyvaluepair and trace parent from env",
			traceParent: tracing.NewTraceParent(),
			args:        []string{"one=true"},
			expectedAttrs: map[string]string{
				"attr.one": "true",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			os.Setenv(TraceParentEnvVar, tc.traceParent)
			if tc.args == nil {
				tc.args = []string{}
			}

			ui := cli.NewMockUi()
			exporter := tracing.NewMemoryExporter()
			cmd, _ := NewAttrCommand(ui)
			cmd.testSpanExporter = exporter

			exitCode := cmd.Run(tc.args)
			assert.Equal(t, tc.expectedExit, exitCode, ui.ErrorWriter.String())

			if tc.expectedError != "" {
				assert.Contains(t, ui.ErrorWriter.String(), tc.expectedError)
			}

			if tc.expectedAttrs != nil {

				state, err := cmd.readState(tc.traceParent)
				assert.NoError(t, err)
				assert.Equal(t, state, tc.expectedAttrs)
			}
		})
	}

}

func TestPairParsing(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		input       []string
		err         string
		result      map[string]string
	}{
		{
			description: "empty input",
			input:       []string{},
			result:      map[string]string{},
		},
		{
			description: "one valid pair",
			input:       []string{"test=yes"},
			result: map[string]string{
				"attr.test": "yes",
			},
		},
		{
			description: "two valid pairs",
			input:       []string{"one=true", "two=false"},
			result: map[string]string{
				"attr.one": "true",
				"attr.two": "false",
			},
		},
		{
			description: "one raw value",
			input:       []string{"one"},
			err:         `'one' was not a valid key=value pair`,
		},
		{
			description: "one key only",
			input:       []string{"one="},
			err:         `'one=' was not a valid key=value pair`,
		},
		{
			description: "one key only",
			input:       []string{"=true"},
			err:         `'=true' was not a valid key=value pair`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {

			result, err := mapFromKeyValues(tc.input)

			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}

			assert.Equal(t, tc.result, result)

		})
	}

}
