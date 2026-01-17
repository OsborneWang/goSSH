package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// TerminalType 表示终端类型
type TerminalType int

const (
	TerminalUnknown TerminalType = iota
	TerminalWindowsTerminal
	TerminalTabby
	TerminalPowerShell
	TerminalCmd
	TerminalITerm2
	TerminalTerminalApp
	TerminalGNOMETerminal
	TerminalKonsole
	TerminalXFCE
	TerminalAlacritty
	TerminalGeneric
)

// String 返回终端类型的字符串表示
func (t TerminalType) String() string {
	switch t {
	case TerminalWindowsTerminal:
		return "Windows Terminal"
	case TerminalTabby:
		return "Tabby"
	case TerminalPowerShell:
		return "PowerShell"
	case TerminalCmd:
		return "CMD"
	case TerminalITerm2:
		return "iTerm2"
	case TerminalTerminalApp:
		return "Terminal.app"
	case TerminalGNOMETerminal:
		return "GNOME Terminal"
	case TerminalKonsole:
		return "Konsole"
	case TerminalXFCE:
		return "XFCE Terminal"
	case TerminalAlacritty:
		return "Alacritty"
	case TerminalGeneric:
		return "Generic"
	default:
		return "Unknown"
	}
}

// DetectTerminal 检测当前终端类型
func DetectTerminal() TerminalType {
	switch runtime.GOOS {
	case "windows":
		return detectWindowsTerminal()
	case "darwin":
		return detectMacOSTerminal()
	case "linux":
		return detectLinuxTerminal()
	default:
		return TerminalUnknown
	}
}

// detectWindowsTerminal 检测 Windows 终端类型
func detectWindowsTerminal() TerminalType {
	// 检查 Windows Terminal (通过 WT_SESSION 环境变量)
	if os.Getenv("WT_SESSION") != "" {
		return TerminalWindowsTerminal
	}

	// 检查 Tabby (通过进程名或环境变量)
	terminal := os.Getenv("TERM_PROGRAM")
	if terminal == "Tabby" || terminal == "tabby" {
		return TerminalTabby
	}

	// 尝试通过父进程检测
	// 这里简化处理，通过环境变量判断
	if os.Getenv("WT_PROFILE_ID") != "" {
		return TerminalWindowsTerminal
	}

	// 检查是否是 PowerShell
	psModulePath := os.Getenv("PSModulePath")
	if psModulePath != "" {
		return TerminalPowerShell
	}

	// 默认认为是 CMD
	return TerminalCmd
}

// detectMacOSTerminal 检测 macOS 终端类型
func detectMacOSTerminal() TerminalType {
	termProgram := os.Getenv("TERM_PROGRAM")
	switch termProgram {
	case "iTerm.app":
		return TerminalITerm2
	case "Apple_Terminal":
		return TerminalTerminalApp
	case "Tabby":
		return TerminalTabby
	default:
		return TerminalTerminalApp // 默认为 Terminal.app
	}
}

// detectLinuxTerminal 检测 Linux 终端类型
func detectLinuxTerminal() TerminalType {
	// 检查环境变量
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	
	if termProgram != "" {
		switch termProgram {
		case "gnome-terminal":
			return TerminalGNOMETerminal
		case "konsole":
			return TerminalKonsole
		}
	}

	// 通过父进程检测（简化版本，实际可能需要更复杂的检测）
	// 检查常见的终端进程
	if strings.Contains(term, "gnome") || strings.Contains(term, "xterm-gnome") {
		return TerminalGNOMETerminal
	}

	// 检查 KDE Konsole
	if strings.Contains(term, "konsole") {
		return TerminalKonsole
	}

	// 检查 XFCE Terminal
	if strings.Contains(term, "xfce") {
		return TerminalXFCE
	}

	// 检查 Alacritty
	if strings.Contains(term, "alacritty") {
		return TerminalAlacritty
	}

	return TerminalGeneric
}

// HasDesktopEnvironment 检查是否有桌面环境（Linux）
func HasDesktopEnvironment() bool {
	if runtime.GOOS != "linux" {
		return true // Windows 和 macOS 默认有桌面环境
	}

	// 检查 DISPLAY (X11)
	if os.Getenv("DISPLAY") != "" {
		return true
	}

	// 检查 WAYLAND_DISPLAY (Wayland)
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}

	// 检查 XDG_CURRENT_DESKTOP
	if os.Getenv("XDG_CURRENT_DESKTOP") != "" {
		return true
	}

	return false
}

// OpenInNewTab 在当前终端的新标签页中执行命令
func OpenInNewTab(command string, args ...string) error {
	terminal := DetectTerminal()
	hasDesktop := HasDesktopEnvironment()

	// 如果没有桌面环境，回退到新窗口
	if !hasDesktop && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		return OpenInNewWindow(command, args...)
	}

	switch terminal {
	case TerminalWindowsTerminal:
		return openInNewTabWindowsTerminal(command, args...)
	case TerminalITerm2:
		return openInNewTabITerm2(command, args...)
	case TerminalTerminalApp:
		return openInNewTabTerminalApp(command, args...)
	case TerminalGNOMETerminal:
		return openInNewTabGNOMETerminal(command, args...)
	case TerminalKonsole:
		return openInNewTabKonsole(command, args...)
	case TerminalTabby:
		// Tabby 可能不支持外部控制新标签页，回退到新窗口
		return OpenInNewWindow(command, args...)
	default:
		// 其他终端尝试新窗口方式
		return OpenInNewWindow(command, args...)
	}
}

// OpenInNewWindow 在新窗口中执行命令
func OpenInNewWindow(command string, args ...string) error {
	terminal := DetectTerminal()
	hasDesktop := HasDesktopEnvironment()

	// 如果没有桌面环境，无法打开新窗口
	if !hasDesktop && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		return fmt.Errorf("无桌面环境，无法打开新窗口")
	}

	switch terminal {
	case TerminalWindowsTerminal:
		return openInNewWindowWindowsTerminal(command, args...)
	case TerminalITerm2:
		return openInNewWindowITerm2(command, args...)
	case TerminalTerminalApp:
		return openInNewWindowTerminalApp(command, args...)
	case TerminalGNOMETerminal:
		return openInNewWindowGNOMETerminal(command, args...)
	case TerminalKonsole:
		return openInNewWindowKonsole(command, args...)
	case TerminalTabby:
		return openInNewWindowTabby(command, args...)
	default:
		// 尝试通用方法
		return openInNewWindowGeneric(command, args...)
	}
}

// openInNewTabWindowsTerminal 在 Windows Terminal 新标签页中执行命令
func openInNewTabWindowsTerminal(command string, args ...string) error {
	// 构建完整命令
	cmdArgs := []string{"-w", "0", "new-tab", command}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("wt.exe", cmdArgs...)
	return cmd.Start() // 不等待，立即返回
}

// openInNewWindowWindowsTerminal 在 Windows Terminal 新窗口中执行命令
func openInNewWindowWindowsTerminal(command string, args ...string) error {
	cmdArgs := []string{"new-window", command}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("wt.exe", cmdArgs...)
	return cmd.Start()
}

// openInNewTabITerm2 在 iTerm2 新标签页中执行命令
func openInNewTabITerm2(command string, args ...string) error {
	// 构建完整命令字符串
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	// 转义命令中的特殊字符
	fullCommand = strings.ReplaceAll(fullCommand, "\\", "\\\\")
	fullCommand = strings.ReplaceAll(fullCommand, "\"", "\\\"")

	// 构建 AppleScript
	script := fmt.Sprintf(`tell application "iTerm2"
	tell current window
		create tab with default profile command "%s"
	end tell
end tell`, fullCommand)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Start()
}

// openInNewWindowITerm2 在 iTerm2 新窗口中执行命令
func openInNewWindowITerm2(command string, args ...string) error {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	fullCommand = strings.ReplaceAll(fullCommand, "\\", "\\\\")
	fullCommand = strings.ReplaceAll(fullCommand, "\"", "\\\"")

	script := fmt.Sprintf(`tell application "iTerm2"
	create window with default profile command "%s"
end tell`, fullCommand)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Start()
}

// openInNewTabTerminalApp 在 Terminal.app 新标签页中执行命令
func openInNewTabTerminalApp(command string, args ...string) error {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	fullCommand = strings.ReplaceAll(fullCommand, "\\", "\\\\")
	fullCommand = strings.ReplaceAll(fullCommand, "\"", "\\\"")

	script := fmt.Sprintf(`tell application "Terminal"
	if (count of windows) = 0 then
		do script "%s"
	else
		tell application "System Events" to keystroke "t" using command down
		do script "%s" in front window
	end if
end tell`, fullCommand, fullCommand)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Start()
}

// openInNewWindowTerminalApp 在 Terminal.app 新窗口中执行命令
func openInNewWindowTerminalApp(command string, args ...string) error {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	fullCommand = strings.ReplaceAll(fullCommand, "\\", "\\\\")
	fullCommand = strings.ReplaceAll(fullCommand, "\"", "\\\"")

	script := fmt.Sprintf(`tell application "Terminal"
	do script "%s"
	activate
end tell`, fullCommand)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Start()
}

// openInNewTabGNOMETerminal 在 GNOME Terminal 新标签页中执行命令
func openInNewTabGNOMETerminal(command string, args ...string) error {
	cmdArgs := []string{"--tab", "--"}
	cmdArgs = append(cmdArgs, command)
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("gnome-terminal", cmdArgs...)
	return cmd.Start()
}

// openInNewWindowGNOMETerminal 在 GNOME Terminal 新窗口中执行命令
func openInNewWindowGNOMETerminal(command string, args ...string) error {
	cmdArgs := []string{"--"}
	cmdArgs = append(cmdArgs, command)
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("gnome-terminal", cmdArgs...)
	return cmd.Start()
}

// openInNewTabKonsole 在 Konsole 新标签页中执行命令
func openInNewTabKonsole(command string, args ...string) error {
	// Konsole 使用 --new-tab 参数
	// 需要构建命令字符串
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	cmdArgs := []string{"--new-tab", "-e", fullCommand}
	cmd := exec.Command("konsole", cmdArgs...)
	return cmd.Start()
}

// openInNewWindowKonsole 在 Konsole 新窗口中执行命令
func openInNewWindowKonsole(command string, args ...string) error {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	cmdArgs := []string{"-e", fullCommand}
	cmd := exec.Command("konsole", cmdArgs...)
	return cmd.Start()
}

// openInNewWindowTabby 在 Tabby 新窗口中执行命令
func openInNewWindowTabby(command string, args ...string) error {
	// Tabby 可能没有直接的 CLI 接口，尝试使用系统默认方式
	// 或者尝试通过 tabby 命令启动
	cmdArgs := []string{"open", "-a", "Tabby", "--args"}
	cmdArgs = append(cmdArgs, command)
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("open", cmdArgs...)
	if err := cmd.Run(); err != nil {
		// 如果失败，尝试通用方法
		return openInNewWindowGeneric(command, args...)
	}
	return nil
}

// openInNewWindowGeneric 使用通用方法在新窗口中执行命令
func openInNewWindowGeneric(command string, args ...string) error {
	// 尝试使用系统默认终端启动器
	switch runtime.GOOS {
	case "windows":
		// Windows: 使用 start 命令
		cmdArgs := []string{"/c", "start", "cmd", "/k"}
		cmdArgs = append(cmdArgs, command)
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.Command("cmd", cmdArgs...)
		return cmd.Start()
	case "darwin":
		// macOS: 使用 open 命令
		fullCommand := command
		if len(args) > 0 {
			fullCommand = command + " " + strings.Join(args, " ")
		}
		script := fmt.Sprintf(`tell application "Terminal"
	do script "%s"
	activate
end tell`, fullCommand)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Start()
	case "linux":
		// Linux: 尝试使用 x-terminal-emulator 或 gnome-terminal
		cmdArgs := []string{"--"}
		cmdArgs = append(cmdArgs, command)
		cmdArgs = append(cmdArgs, args...)

		// 尝试 gnome-terminal
		if _, err := exec.LookPath("gnome-terminal"); err == nil {
			cmd := exec.Command("gnome-terminal", cmdArgs...)
			return cmd.Start()
		}

		// 尝试 x-terminal-emulator
		if _, err := exec.LookPath("x-terminal-emulator"); err == nil {
			cmd := exec.Command("x-terminal-emulator", cmdArgs...)
			return cmd.Start()
		}

		return fmt.Errorf("未找到可用的终端模拟器")
	default:
		return fmt.Errorf("不支持的操作系统")
	}
}

// GetExecutablePath 获取可执行文件路径
func GetExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %v", err)
	}

	return absPath, nil
}
