package utils

import (
	"os/exec"
	"runtime"
)

var (
	execCommand = exec.Command
	runtimeOS   = runtime.GOOS
)

func OpenURL(url string) error {
	var cmd string
	switch runtimeOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return execCommand(cmd, url).Start()
}
