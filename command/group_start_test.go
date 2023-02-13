package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestSpanStartArgumentParsing(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description   string
		traceParent   string
		args          []string
		expectedExit  int
		expectedError string
	}{
		{
			description:   "no args and no environment",
			expectedExit:  1,
			expectedError: "this command takes at least 1 argument: group_name",
		},
		{
			description:   "only span name, no trace from env",
			args:          []string{"test_span"},
			expectedExit:  1,
			expectedError: "this command requires a trace_parent",
		},
		{
			description:  "span name and trace parent",
			args:         []string{"test_span", NewTraceParent()},
			expectedExit: 0,
		},
		{
			description:  "span name and trace parent from env",
			traceParent:  NewTraceParent(),
			args:         []string{"test_span"},
			expectedExit: 0,
		},
		{
			description:   "traceparent is invalid format",
			args:          []string{"test_span"},
			traceParent:   "00-00000000000000000000000000000000-0000000000000000-00",
			expectedExit:  1,
			expectedError: "trace-id can't be all zero",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			os.Setenv(TraceParentEnvVar, tc.traceParent)
			if tc.args == nil {
				tc.args = []string{}
			}

			ui := cli.NewMockUi()
			cmd, _ := NewGroupStartCommand(ui)

			exitCode := cmd.Run(tc.args)
			assert.Equal(t, tc.expectedExit, exitCode, ui.ErrorWriter.String())

			if tc.expectedError != "" {
				assert.Contains(t, ui.ErrorWriter.String(), tc.expectedError)
			}
		})
	}
}

func TestSpanStart(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	cmd, _ := NewGroupStartCommand(ui)

	now := time.Now().UnixNano()
	cmd.Base.now = func() int64 {
		return now
	}

	parentTrace := NewTraceID()
	parentSpan := NewSpanID()

	assert.Equal(t, 0, cmd.Run([]string{"test-generate", AsTraceParent(parentTrace, parentSpan)}))

	traceParent := ui.OutputWriter.String()
	filepath := path.Join(os.TempDir(), "trace", "state", strings.TrimSpace(traceParent))

	content, err := ioutil.ReadFile(filepath)
	assert.NoError(t, err)

	assert.Contains(t, string(content), "name=test-generate")
	assert.Contains(t, string(content), fmt.Sprintf("start=%v", now))
	assert.Contains(t, string(content), fmt.Sprintf("parent=%s", parentSpan))
}
