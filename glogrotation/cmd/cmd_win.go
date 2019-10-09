//+build windows

package cmd

import (
	"os"
	"os/exec"
)

func prepareCmd(cmd *exec.Cmd) {

}

func killable(sig os.Signal) bool {
	return true
}
