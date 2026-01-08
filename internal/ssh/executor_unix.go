// +build !windows

package ssh

import (
	"os"

	"golang.org/x/sys/unix"
)

// getTerminalSize 获取终端大小（Unix系统）
func getTerminalSize(fd int) (width, height int) {
	var ws unix.Winsize
	if err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ); err == nil {
		width = int(ws.Col)
		height = int(ws.Row)
	}
	return width, height
}


