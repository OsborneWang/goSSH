# 开发指南

本文档介绍如何在开发 GoSSH 项目时进行调试和测试。

## 🚀 快速开始开发

### 方法 1: 直接运行（推荐用于快速测试）

不需要每次修改都重新编译，直接使用 `go run` 命令：

```bash
# 运行主程序（显示帮助）
go run main.go

# 运行交互式模式
go run main.go interactive

# 运行其他命令
go run main.go list
go run main.go add
go run main.go connect server1
go run main.go exec server1 "ls -la"
```

**优点：** 
- 不需要编译，快速看到结果
- 自动处理依赖

**缺点：**
- 每次运行都会重新编译（但Go的增量编译很快）

### 方法 2: 使用 Makefile（推荐）

项目提供了 Makefile，包含常用的开发命令：

```bash
# 直接运行（开发模式）
make run

# 运行交互式模式
make run-interactive

# 运行列表命令
make run-list

# 构建可执行文件
make build

# 格式化代码
make fmt

# 代码检查
make vet

# 清理构建文件
make clean

# 查看所有可用命令
make help
```

### 方法 3: 热重载开发（推荐用于频繁修改）

使用 [Air](https://github.com/cosmtrek/air) 工具，代码修改后自动重新编译运行：

```bash
# 1. 首先安装 Air
go install github.com/cosmtrek/air@latest

# 或者使用 Makefile 安装开发工具
make install-tools

# 2. 启动热重载开发模式
make watch
# 或直接运行
air

# 3. 修改代码后，Air 会自动检测并重新编译运行
```

Air 配置文件已包含在项目中（`.air.toml`），它会监听 `.go` 文件的改动并自动重启。

## 🐛 调试方法

### 方法 1: 使用 VS Code 调试器（推荐）

1. **安装 VS Code Go 扩展**
   - 安装官方 Go 扩展：`ms-vscode.go`

2. **配置调试**
   - 项目已包含 `.vscode/launch.json` 调试配置
   - 按 `F5` 开始调试，或在调试面板选择配置

3. **可用的调试配置：**
   - `Launch GoSSH` - 运行主程序（显示帮助）
   - `Launch GoSSH (Interactive Mode)` - 运行交互式模式
   - `Launch GoSSH (Connect)` - 测试连接功能
   - `Launch GoSSH (Exec)` - 测试执行命令功能

4. **设置断点：**
   - 在代码行号左侧点击设置断点
   - 使用 `F9` 切换断点
   - 使用 `F10` 单步跳过，`F11` 单步进入

5. **调试变量：**
   - 在左侧面板查看变量值
   - 在"调试控制台"中输入变量名查看值
   - 使用鼠标悬停查看变量值

### 方法 2: 使用命令行调试器 Delve

```bash
# 1. 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 2. 使用 dlv 启动调试
dlv debug .

# 3. 在 dlv 命令行中：
# (dlv) break main.main          # 在 main 函数设置断点
# (dlv) break cmd/connect.go:30  # 在指定位置设置断点
# (dlv) continue                 # 继续执行
# (dlv) next                     # 下一行
# (dlv) step                     # 进入函数
# (dlv) print variable_name      # 打印变量
# (dlv) exit                     # 退出

# 或者直接调试特定命令
dlv debug . -- interactive
dlv debug . -- connect server1
```

### 方法 3: 使用日志输出调试

在代码中添加日志输出：

```go
import (
    "fmt"
    "log"
)

// 使用标准库
fmt.Printf("调试: 变量值 = %v\n", variable)
log.Printf("调试: 执行到这里了\n")

// 使用详细日志
log.Printf("[DEBUG] Server: %+v\n", server)
log.Printf("[DEBUG] Connection status: %v\n", client.IsConnected())
```

可以在代码中临时添加这些日志，调试完成后删除。

### 方法 4: 使用环境变量控制调试

在代码中添加调试标志：

```go
package main

import (
    "os"
    "log"
)

var debug = os.Getenv("GOSSH_DEBUG") == "1"

func debugLog(format string, v ...interface{}) {
    if debug {
        log.Printf("[DEBUG] "+format, v...)
    }
}

// 使用时
debugLog("连接服务器: %s", server.Name)
```

运行时启用调试：
```bash
GOSSH_DEBUG=1 go run main.go connect server1
```

**Windows PowerShell:**
```powershell
$env:GOSSH_DEBUG="1"; go run main.go connect server1
```

**Windows CMD:**
```cmd
set GOSSH_DEBUG=1 && go run main.go connect server1
```

## 📝 开发工作流

### 1. 日常开发流程

```bash
# 1. 启动热重载（推荐）
make watch

# 或者在另一个终端运行
air

# 2. 在代码编辑器中修改代码
# 3. 保存文件，Air 会自动重新编译运行
# 4. 测试功能
```

### 2. 测试新功能

```bash
# 1. 修改代码后，直接运行测试
go run main.go <command> <args>

# 2. 或者先编译再运行（用于性能测试）
make build
./goss.exe <command> <args>
```

### 3. 代码质量检查

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 运行 linter（需要先安装 golangci-lint）
make lint

# 运行测试（如果有测试文件）
make test
```

## 🔧 开发工具安装

### 必需工具

```bash
# Go 语言（1.21+）
go version

# 安装项目依赖
go mod download
```

### 推荐工具

```bash
# 安装所有开发工具
make install-tools

# 或手动安装：

# 1. goimports - 自动管理导入
go install golang.org/x/tools/cmd/goimports@latest

# 2. golangci-lint - 代码检查工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 3. Air - 热重载开发
go install github.com/cosmtrek/air@latest

# 4. Delve - 调试器
go install github.com/go-delve/delve/cmd/dlv@latest
```

## 📂 项目结构说明

```
goSSH/
├── main.go                 # 程序入口
├── cmd/                    # Cobra命令定义
│   ├── root.go            # 根命令
│   ├── add.go             # 添加服务器
│   ├── list.go            # 列出服务器
│   ├── connect.go         # 连接服务器
│   ├── exec.go            # 执行命令
│   ├── transfer.go        # 文件传输
│   └── interactive.go     # 交互式模式
├── internal/
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── ssh/               # SSH功能
│   │   ├── client.go      # SSH客户端
│   │   ├── executor.go    # 命令执行
│   │   ├── executor_unix.go  # Unix终端大小
│   │   └── transfer.go    # 文件传输
│   └── storage/           # 存储
│       └── storage.go
├── models/                # 数据模型
│   └── server.go
├── .vscode/               # VS Code配置
│   ├── launch.json        # 调试配置
│   └── settings.json      # 编辑器设置
├── Makefile               # 构建脚本
├── .air.toml             # Air配置
└── README.md             # 使用文档
```

## 🎯 调试技巧

### 1. 调试 SSH 连接问题

```go
// 在 internal/ssh/client.go 的 Connect 方法中添加
func (c *Client) Connect() error {
    fmt.Printf("[DEBUG] 连接到 %s:%d\n", c.server.Host, c.server.Port)
    fmt.Printf("[DEBUG] 用户: %s\n", c.server.Username)
    
    // ... 连接代码
    
    if err != nil {
        fmt.Printf("[DEBUG] 连接错误详情: %+v\n", err)
        return err
    }
    
    fmt.Printf("[DEBUG] 连接成功\n")
    return nil
}
```

### 2. 调试配置文件读取

```go
// 在 internal/storage/storage.go 中添加
func (s *Storage) Load() (*models.ServerConfig, error) {
    fmt.Printf("[DEBUG] 配置文件路径: %s\n", s.configPath)
    
    // ... 加载代码
    
    fmt.Printf("[DEBUG] 加载了 %d 个服务器\n", len(config.Servers))
    return config, nil
}
```

### 3. 调试命令执行

```go
// 在 internal/ssh/executor.go 中添加
func (e *Executor) Execute(command string) (string, error) {
    fmt.Printf("[DEBUG] 执行命令: %s\n", command)
    
    output, err := session.CombinedOutput(command)
    
    fmt.Printf("[DEBUG] 输出长度: %d 字节\n", len(output))
    fmt.Printf("[DEBUG] 错误: %v\n", err)
    
    return string(output), err
}
```

### 4. 使用条件编译调试

在调试版本中添加更多信息：

```go
//go:build debug
// +build debug

package ssh

const debugMode = true
```

编译调试版本：
```bash
go build -tags debug -o goss-debug.exe .
```

## 🧪 测试服务器

开发时需要测试连接功能，可以：

1. **使用本地SSH服务器**（需要安装OpenSSH）
   ```bash
   # Windows 10+
   # 在"设置" -> "应用" -> "可选功能" 中安装 OpenSSH 服务器
   
   # Linux
   sudo apt install openssh-server
   sudo systemctl start sshd
   
   # 添加测试服务器
   goss add
   # 输入: localhost, 127.0.0.1, 22, your_username, your_password
   ```

2. **使用Docker容器作为测试服务器**
   ```bash
   docker run -d -p 2222:22 --name test-ssh \
     -e ROOT_PASSWORD=testpass123 \
     panubo/sshd
   
   # 添加服务器
   goss add
   # 输入: docker-test, localhost, 2222, root, testpass123
   ```

3. **使用云服务器或虚拟机**

## 💡 开发建议

1. **频繁提交代码** - 使用 Git 管理代码，频繁提交
2. **测试每个功能** - 开发完一个功能立即测试
3. **查看日志** - 注意程序的输出信息
4. **使用断点** - 在关键位置设置断点进行调试
5. **代码格式化** - 提交前运行 `make fmt`
6. **代码检查** - 运行 `make vet` 检查潜在问题

## 🔍 常见调试场景

### 场景 1: 修改了代码但没看到效果

**解决方法：**
```bash
# 确保重新编译了
make clean
make build

# 或使用 go run 确保使用最新代码
go run main.go
```

### 场景 2: 配置文件问题

```bash
# 检查配置文件位置
echo %APPDATA%\gossh\servers.json  # Windows
echo ~/.config/gossh/servers.json  # Linux/macOS

# 查看配置文件内容
cat ~/.config/gossh/servers.json   # Linux/macOS
type %APPDATA%\gossh\servers.json  # Windows
```

### 场景 3: 连接超时或失败

```bash
# 使用调试模式运行
GOSSH_DEBUG=1 go run main.go connect server1

# 检查网络连接
ping server_host
telnet server_host 22
```

### 场景 4: 终端大小检测问题

Windows 和 Unix 系统的终端大小获取方式不同，已使用条件编译处理。如果遇到问题，检查 `executor.go` 和 `executor_unix.go`。

---

**祝开发愉快！** 🎉
