package metadata

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestParseCommandLine(t *testing.T) {
	assert := require.New(t)

	tests := []struct {
		In    string
		Out   []string
		Error error
	}{
		{
			In:  "/path/to/bin arg1 arg2 arg3",
			Out: []string{"/path/to/bin", "arg1", "arg2", "arg3"},
		},
	}

	for _, tst := range tests {
		out, err := parseCommandLine(tst.In)
		assert.Equal(tst.Error, err)
		if err != nil {
			continue
		}
		assert.Equal(tst.Out, out)
	}
}

func TestRunCommand(t *testing.T) {
	assert := require.New(t)

	tests := []struct {
		In           string
		Out          string
		MaxExecution time.Duration
		Error        error
	}{
		{
			In:           "echo foo bar",
			Out:          "foo bar",
			MaxExecution: time.Second,
		},
		{
			In:           "sleep 2",
			MaxExecution: time.Second,
			Error:        errors.New("execution error: signal: killed"),
		},
	}

	for _, tst := range tests {
		maxExecution = tst.MaxExecution
		out, err := runCommand(tst.In)
		if err != nil || tst.Error != nil {
			assert.Equal(tst.Error.Error(), err.Error())
		}
		if err != nil {
			continue
		}
		assert.Equal(tst.Out, out)
	}
}

func TestMetaData(t *testing.T) {
	tests := []struct {
		Name     string
		Static   map[string]string
		Commands map[string]string
		Expected map[string]string
	}{
		{
			Name: "static only",
			Static: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
			Expected: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
		},
		{
			Name: "commands only",
			Commands: map[string]string{
				"foo": "echo test1",
				"bar": "echo test2",
			},
			Expected: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
		},
		{
			Name: "static + commands",
			Static: map[string]string{
				"static_1": "static 1",
				"static_2": "static_2",
			},
			Commands: map[string]string{
				"cmd_1": "echo cmd1",
				"cmd_2": "echo cmd2",
			},
			Expected: map[string]string{
				"static_1": "static 1",
				"static_2": "static_2",
				"cmd_1":    "cmd1",
				"cmd_2":    "cmd2",
			},
		},
		{
			Name: "command overwrites static",
			Static: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
			Commands: map[string]string{
				"bar": "echo cmd overwrite",
			},
			Expected: map[string]string{
				"foo": "test1",
				"bar": "cmd overwrite",
			},
		},
		{
			Name: "command overwrites but timeout",
			Static: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
			Commands: map[string]string{
				"bar": "sleep 2",
			},
			Expected: map[string]string{
				"foo": "test1",
				"bar": "test2",
			},
		},
	}

	maxExecution = time.Second

	for _, tst := range tests {
		t.Run(tst.Name, func(t *testing.T) {
			assert := require.New(t)

			static = tst.Static
			commands = tst.Commands

			runCommands()

			assert.EqualValues(tst.Expected, Get())
		})
	}
}
