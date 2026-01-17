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

// crlfFilterReader 过滤掉Windows终端发送的\r字符，只保留\n
// 这样可以避免在SSH会话中出现双重回车的问题
type crlfFilterReader struct {
	reader io.Reader
}

func (r *crlfFilterReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		// 过滤掉\r字符
		writeIdx := 0
		for i := 0; i < n; i++ {
			if p[i] != '\r' {
				p[writeIdx] = p[i]
				writeIdx++
			}
		}
		n = writeIdx
	}
	return n, err
}

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
	// 在Windows上，使用crlfFilterReader过滤掉\r字符，避免双重回车问题
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = &crlfFilterReader{reader: os.Stdin}

	// 设置伪终端（PTY）用于交互式命令
	// ECHO 设置为 0，禁用远程回显，由本地终端负责回显，避免命令重复显示
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用远程回显，避免命令重复显示
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

// executeShellInCurrentTerminal 在当前终端中启动交互式Shell（原有逻辑）
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

	// 设置标准输入输出
	// 在Windows上，使用crlfFilterReader过滤掉\r字符，避免双重回车问题
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = &crlfFilterReader{reader: os.Stdin}

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
