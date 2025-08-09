package utils

import (
	"os/exec"
	"runtime"
)

func OpenURI(uri string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", uri}
	case "darwin":
		cmd = "open"
		args = []string{uri}
	default: // Linux and other Unix-like system.
		cmd = "xdg-open"
		args = []string{uri}
	}
	return exec.Command(cmd, args...).Start()
}
