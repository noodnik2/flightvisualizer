//go:build darwin

package os

import "os/exec"

func LaunchFile(filename string) error {
	return exec.Command("open", filename).Run()
}
