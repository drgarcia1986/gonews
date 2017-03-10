package utils

import (
	"os/exec"
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
	f := new(FakeExecCommand)

	execCommand = f.execCommand
	defer func() { execCommand = exec.Command }()

	url := "http://google.com"
	if err := OpenURL(url); err != nil {
		t.Fatal("error on exec fake command: ", err)
	}

	if f.args[0] != url {
		t.Errorf("expected %s, got %s", url, f.args[0])
	}
}
