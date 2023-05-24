package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestStarting(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	cmd, _ := NewStartCommand(ui)

	now := time.Now().UnixNano()
	cmd.Base.now = func() int64 {
		return now
	}

	assert.Equal(t, 0, cmd.Run([]string{"test-generate", "--attr", "branch=testing", "--attr", "trigger=cron"}))

	traceParent := ui.OutputWriter.String()
	assert.Regexp(t, `^[[:xdigit:]]{2}-[[:xdigit:]]{32}-[[:xdigit:]]{16}-[[:xdigit:]]{2}\n`, traceParent)

	filepath := path.Join(os.TempDir(), "trace", "state", strings.TrimSpace(traceParent))

	content, err := ioutil.ReadFile(filepath)
	assert.NoError(t, err)

	assert.Contains(t, string(content), "name=test-generate")
	assert.Contains(t, string(content), fmt.Sprintf("start=%v", now))
	assert.Contains(t, string(content), "attr.branch=testing")
	assert.Contains(t, string(content), "attr.trigger=cron")
}

func TestTimestamps(t *testing.T) {
	t.Parallel()

	when := time.Now().Add(-30 * time.Second)

	t.Run("iso-8601", func(t *testing.T) {
		ui := cli.NewMockUi()
		cmd, _ := NewStartCommand(ui)

		now := time.Now().UnixNano()
		cmd.Base.now = func() int64 {
			return now
		}

		exitCode := cmd.Run([]string{"custom-time-iso", "--when", when.Format(ISO8601)})
		assert.Zero(t, exitCode, "should not error when running")
		traceParent := strings.TrimSpace(ui.OutputWriter.String())

		state, err := cmd.readState(traceParent)
		assert.NoError(t, err)

		assert.Equal(t, strconv.FormatInt(when.Truncate(time.Second).UnixNano(), 10), state["start"])
	})

	t.Run("epoch", func(t *testing.T) {
		ui := cli.NewMockUi()
		cmd, _ := NewStartCommand(ui)

		now := time.Now().UnixNano()
		cmd.Base.now = func() int64 {
			return now
		}

		exitCode := cmd.Run([]string{"custom-time-epoch", "--when", strconv.FormatInt(when.Unix(), 10)})
		assert.Zero(t, exitCode, "should not error when running")
		traceParent := strings.TrimSpace(ui.OutputWriter.String())

		state, err := cmd.readState(traceParent)
		assert.NoError(t, err)

		assert.Equal(t, strconv.FormatInt(when.Truncate(time.Second).UnixNano(), 10), state["start"])
	})
}

func TestTime(t *testing.T) {

	cases := []struct {
		input  string
		format string
	}{
		{
			input:  "2023-05-22T09:46:00Z",
			format: "RFC3339",
		},
		{
			input:  strconv.FormatInt(time.Now().Add(-30*time.Second).Unix(), 10),
			format: "unix",
		},
		{
			input:  time.Now().Add(-30 * time.Second).Format(ISO8601),
			format: "iso-8601",
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			cmd, _ := NewStartCommand(cli.NewMockUi())
			cmd.startTime = tc.input
			_, err := cmd.startEpochNano()

			// _, err := time.Parse(tc.format, tc.input)
			assert.NoError(t, err)
		})
	}
}
