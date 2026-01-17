// +build !windows

package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
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

// ExecuteShell 启动交互式Shell
// 如果 useNewTab 为 false，则优先尝试在新标签页中启动，失败后尝试新窗口，最后回退到当前终端
func (e *Executor) ExecuteShell(useNewTab bool) error {
	// 如果不需要新标签页，直接在当前终端执行
	if !useNewTab {
		return e.executeShellInCurrentTerminal()
	}

	// 优先尝试新标签页
	if err := e.ExecuteShellInNewTab(); err == nil {
		return nil
	}

	// 如果新标签页失败，尝试新窗口
	if err := e.ExecuteShellInNewWindow(); err == nil {
		return nil
	}

	// 如果都失败，回退到当前终端
	return e.executeShellInCurrentTerminal()
}

// ExecuteShellInNewTab 在新标签页中启动交互式Shell
func (e *Executor) ExecuteShellInNewTab() error {
	// 获取服务器名称
	serverName := e.client.GetServer().Name

	// 获取可执行文件路径
	execPath, err := GetExecutablePath()
	if err != nil {
		return err
	}

	// 构建命令：goss connect servername --no-new-tab
	// --no-new-tab 标志确保在新标签页中不会再尝试打开新标签页
	cmdArgs := []string{"connect", serverName, "--no-new-tab"}

	// 在新标签页中执行命令
	return OpenInNewTab(execPath, cmdArgs...)
}

// ExecuteShellInNewWindow 在新窗口中启动交互式Shell
func (e *Executor) ExecuteShellInNewWindow() error {
	serverName := e.client.GetServer().Name

	execPath, err := GetExecutablePath()
	if err != nil {
		return err
	}

	cmdArgs := []string{"connect", serverName, "--no-new-tab"}

	return OpenInNewWindow(execPath, cmdArgs...)
}

// executeShellInCurrentTerminal 在当前终端中启动交互式Shell（Unix版本，不使用crlfFilterReader）
func (e *Executor) executeShellInCurrentTerminal() error {
	if !e.client.IsConnected() {
		if err := e.client.Connect(); err != nil {
			return err
		}
	}

	session, err := e.client.GetConnection().NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %v", err)
	}
	defer session.Close()

	// 设置标准输入输出（Unix不需要crlfFilterReader）
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 设置伪终端
	// ECHO 设置为 0，禁用远程回显，由本地终端负责回显，避免命令重复显示
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用远程回显，避免命令重复显示
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// 获取终端大小
	fd := int(os.Stdin.Fd())
	width, height := getTerminalSize(fd)
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return fmt.Errorf("请求PTY失败: %v", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("启动Shell失败: %v", err)
	}

	// 等待会话结束
	return session.Wait()
}

