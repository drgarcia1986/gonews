package utils

import (
	"os/exec"
	"runtime"
	"testing"
)

type FakeExecCommand struct {
	cmd  string
	args []string
}

func (f *FakeExecCommand) execCommand(cmd string, args ...string) *exec.Cmd {
	f.cmd = cmd
	f.args = args

	return exec.Command("echo", "", ">", "/dev/null")
}

func TestOpen(t *testing.T) {
	tt := []struct {
		runtime  string
		expected string
	}{
		{"darwin", "open"},
		{"linux", "xdg-open"},
	}

	f := new(FakeExecCommand)

	execCommand = f.execCommand
	defer func() { execCommand = exec.Command }()
	defer func() { runtimeOS = runtime.GOOS }()

	url := "http://google.com"
	for _, tc := range tt {
		runtimeOS = tc.runtime
		if err := OpenURL(url); err != nil {
			t.Fatal("error on exec fake command: ", err)
		}

		if f.cmd != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, f.cmd)
		}

		if f.args[0] != url {
			t.Errorf("expected %s, got %s", url, f.args[0])
		}

	}
}
