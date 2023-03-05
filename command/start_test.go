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
