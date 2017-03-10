package utils

import (
	"os/exec"
	"runtime"
)

var execCommand = exec.Command

func OpenURL(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return execCommand(cmd, url).Start()
}
