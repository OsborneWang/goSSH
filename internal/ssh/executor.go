// +build windows

package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
)

// Executor 提供远程命令执行功能
type Executor struct {
	client *Client
}

// NewExecutor 创建新的命令执行器
func NewExecutor(client *Client) *Executor {
	return &Executor{client: client}
}

// Execute 执行远程命令并返回输出
func (e *Executor) Execute(command string) (string, error) {
	if !e.client.IsConnected() {
		if err := e.client.Connect(); err != nil {
			return "", err
		}
	}

	session, err := e.client.GetConnection().NewSession()
	if err != nil {
		return "", fmt.Errorf("创建会话失败: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("执行命令失败: %v", err)
	}

	return string(output), nil
}

// ExecuteInteractive 交互式执行命令（实时输出）
func (e *Executor) ExecuteInteractive(command string) error {
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

	// 设置标准输入输出
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 设置伪终端（PTY）用于交互式命令
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	// 获取终端大小
	fd := int(os.Stdin.Fd())
	if width, height := getTerminalSize(fd); width > 0 && height > 0 {
		if err := session.RequestPty("xterm", height, width, modes); err != nil {
			return fmt.Errorf("请求PTY失败: %v", err)
		}
	}

	if err := session.Run(command); err != nil {
		return fmt.Errorf("执行命令失败: %v", err)
	}

	return nil
}

// ExecuteWithStream 执行命令并实时流式输出
func (e *Executor) ExecuteWithStream(command string) error {
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

	// 获取标准输出和错误输出
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取标准输出管道失败: %v", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取错误输出管道失败: %v", err)
	}

	// 启动命令
	if err := session.Start(command); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	// 实时输出
	done := make(chan bool, 2)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
		done <- true
	}()

	// 等待输出完成
	<-done
	<-done

	// 等待命令执行完成
	if err := session.Wait(); err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("命令执行失败，退出码: %d", exitErr.ExitStatus())
		}
		return fmt.Errorf("等待命令完成失败: %v", err)
	}

	return nil
}

// ExecuteShell 启动交互式Shell
func (e *Executor) ExecuteShell() error {
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

	// 设置标准输入输出
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 设置伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
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

// CopyOutput 复制输出到指定的writer
func (e *Executor) CopyOutput(command string, writer io.Writer) error {
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

	output, err := session.CombinedOutput(command)
	if err != nil {
		if _, writeErr := writer.Write(output); writeErr != nil {
			return writeErr
		}
		return err
	}

	_, err = writer.Write(output)
	return err
}

// getTerminalSize 获取终端大小（Windows系统）
func getTerminalSize(fd int) (width, height int) {
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err == nil {
		width = int(info.Window.Right - info.Window.Left + 1)
		height = int(info.Window.Bottom - info.Window.Top + 1)
	}
	return width, height
}
