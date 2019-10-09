//+build !windows

package cmd

import (
	"os"
	"os/exec"
	"syscall"
)

func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killable(sig os.Signal) bool {
	if sig != syscall.SIGCHLD {
		return true
	}
	return false
}
