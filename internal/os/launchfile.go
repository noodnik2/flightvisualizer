package os

import (
    "os/exec"
    "runtime"
)

func LaunchFile(filename string) error {
    switch runtime.GOOS {
    case "windows":
        return exec.Command("cmd", "/C", filename).Run()
    case "darwin":
        return exec.Command("open", filename).Run()
    case "linux":
        return exec.Command("sh", "-c", filename).Run()
    }
    return exec.Command(filename).Run()
}
