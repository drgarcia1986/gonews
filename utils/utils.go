package utils

import (
	"os/exec"
	"runtime"
)

const Version = "0.0.1"

func OpenURL(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}
