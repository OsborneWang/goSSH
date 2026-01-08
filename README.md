# GoSSH - 跨平台SSH命令行工具

GoSSH 是一个使用 Go 语言开发的跨平台 SSH 命令行工具，支持 Windows、Linux 和 macOS。它提供了服务器管理、SSH 连接、远程命令执行和文件传输功能，让 SSH 操作更加简单便捷。

## ✨ 功能特性

- 🔐 **服务器管理** - 添加、删除、列出服务器配置
- 🖥️ **SSH 连接** - 支持交互式 Shell 连接
- ⚡ **命令执行** - 在远程服务器上执行命令并实时查看输出
- 📁 **文件传输** - 支持 SFTP 上传/下载文件和目录
- 🎯 **交互式模式** - 友好的交互式菜单界面
- 🌐 **跨平台支持** - 支持 Windows、Linux、macOS

## 📦 安装

### 从源码编译

1. 确保已安装 Go 1.21 或更高版本
2. 克隆或下载项目到本地
3. 进入项目目录并编译：

```bash
go build -o goss .
```

Windows 系统会生成 `goss.exe`，Linux/macOS 会生成 `goss` 可执行文件。

### 添加到 PATH（可选）

将编译好的可执行文件添加到系统 PATH 中，这样就可以在任何地方使用 `goss` 命令。

**Windows:**
```powershell
# 将 goss.exe 移动到系统 PATH 中的某个目录
# 或创建一个符号链接
```

**Linux/macOS:**
```bash
sudo mv goss /usr/local/bin/
# 或者
sudo cp goss /usr/local/bin/
```

## 🚀 快速开始

### 1. 添加服务器配置

首次使用时，需要添加服务器配置：

```bash
goss add
```

然后按提示输入：
- 服务器名称/别名（用于标识服务器）
- 主机地址（IP 或域名）
- SSH 端口（默认 22）
- 用户名
- 密码

### 2. 查看服务器列表

```bash
goss list
```

### 3. 连接到服务器

```bash
# 交互式选择服务器
goss connect

# 或直接指定服务器名称
goss connect server1
```

### 4. 执行远程命令

```bash
# 交互式选择服务器和输入命令
goss exec

# 指定服务器和命令
goss exec server1 "ls -la"

# 执行多个参数的命令
goss exec server1 "df -h"
```

### 5. 文件传输

**上传文件/目录：**
```bash
# 交互式模式
goss transfer upload

# 直接指定参数
goss transfer upload server1 /local/path/file.txt /remote/path/

# 上传整个目录
goss transfer upload server1 ./local_dir /remote/dir
```

**下载文件/目录：**
```bash
# 交互式模式
goss transfer download

# 直接指定参数
goss transfer download server1 /remote/path/file.txt ./local/

# 下载整个目录
goss transfer download server1 /remote/dir ./local_dir
```

### 6. 交互式菜单模式

进入交互式菜单，可以更方便地使用所有功能：

```bash
goss interactive
```

交互式菜单提供了以下选项：
- 连接服务器 (SSH Shell)
- 执行远程命令
- 上传文件/目录
- 下载文件/目录
- 添加服务器
- 删除服务器
- 列出所有服务器
- 退出

## 📖 命令详解

### `goss add`

交互式添加新的 SSH 服务器配置。

**使用示例：**
```bash
goss add
```

### `goss list`

列出所有已配置的服务器，显示服务器名称、主机地址、端口和用户名。

**使用示例：**
```bash
goss list
```

**输出示例：**
```
名称                 主机                  端口     用户名         
────────────────────────────────────────────────────────────
server1              192.168.1.100         22       root          
server2              example.com           2222     admin         
```

### `goss remove [name]`

删除指定的服务器配置。如果不提供名称，会进入交互式选择。

**使用示例：**
```bash
# 交互式选择
goss remove

# 直接指定名称
goss remove server1
```

### `goss connect [name]`

连接到 SSH 服务器并启动交互式 Shell。如果不提供名称，会进入交互式选择。

**使用示例：**
```bash
# 交互式选择
goss connect

# 直接指定名称
goss connect server1
```

连接成功后，您将进入远程服务器的 Shell，可以执行各种命令。输入 `exit` 或按 `Ctrl+D` 退出。

### `goss exec [name] [command]`

在远程服务器上执行命令并实时显示输出。

**使用示例：**
```bash
# 交互式模式
goss exec

# 指定服务器和命令
goss exec server1 "ls -la"

# 执行系统命令
goss exec server1 "df -h"
goss exec server1 "ps aux"
```

### `goss transfer upload [name] [local] [remote]`

上传本地文件或目录到远程服务器。

**使用示例：**
```bash
# 交互式模式
goss transfer upload

# 上传文件
goss transfer upload server1 ./local.txt /home/user/

# 上传目录（会递归上传所有文件）
goss transfer upload server1 ./local_dir /home/user/remote_dir
```

### `goss transfer download [name] [remote] [local]`

从远程服务器下载文件或目录到本地。

**使用示例：**
```bash
# 交互式模式
goss transfer download

# 下载文件
goss transfer download server1 /home/user/remote.txt ./

# 下载目录（会递归下载所有文件）
goss transfer download server1 /home/user/remote_dir ./local_dir
```

### `goss interactive`

进入交互式菜单模式，提供友好的菜单界面来执行各种操作。

**使用示例：**
```bash
goss interactive
```

## 📁 配置文件

服务器配置以 JSON 格式存储在以下位置：

- **Windows:** `%APPDATA%\gossh\servers.json`
- **Linux/macOS:** `~/.config/gossh/servers.json`

配置文件格式示例：

```json
{
  "servers": [
    {
      "name": "server1",
      "host": "192.168.1.100",
      "port": 22,
      "username": "root",
      "password": "your_password"
    },
    {
      "name": "server2",
      "host": "example.com",
      "port": 2222,
      "username": "admin",
      "password": "another_password"
    }
  ]
}
```

⚠️ **安全提示：** 密码以明文形式存储。请确保配置文件权限设置正确，不要在公共环境中使用此工具存储敏感服务器信息。

## 🔧 技术栈

- **Go 1.21+**
- **golang.org/x/crypto/ssh** - SSH 客户端库
- **github.com/pkg/sftp** - SFTP 文件传输
- **github.com/spf13/cobra** - 命令行框架
- **github.com/manifoldco/promptui** - 交互式输入
- **github.com/fatih/color** - 终端颜色输出

## 📝 使用示例

### 示例 1：完整的服务器管理流程

```bash
# 1. 添加服务器
goss add
# 输入: server1, 192.168.1.100, 22, root, password123

# 2. 查看服务器列表
goss list

# 3. 连接服务器
goss connect server1

# 4. 执行命令
goss exec server1 "uptime"

# 5. 上传文件
goss transfer upload server1 ./script.sh /tmp/

# 6. 下载文件
goss transfer download server1 /var/log/app.log ./
```

### 示例 2：使用交互式模式

```bash
# 进入交互式菜单
goss interactive

# 在菜单中选择操作
# 1. 连接服务器 (SSH Shell)
# 2. 执行远程命令
# 3. 上传文件/目录
# 4. 下载文件/目录
# ...
```

### 示例 3：批量操作

```bash
# 在多个服务器上执行相同命令
goss exec server1 "sudo apt update"
goss exec server2 "sudo apt update"
goss exec server3 "sudo apt update"
```

## ⚠️ 注意事项

1. **安全性：** 
   - 密码以明文形式存储在配置文件中
   - 建议在生产环境中使用 SSH 密钥认证（未来版本可能支持）
   - 确保配置文件权限设置正确

2. **网络连接：**
   - 确保网络连接正常
   - 确保防火墙允许 SSH 连接
   - 检查服务器 SSH 服务是否运行

3. **文件传输：**
   - 大文件传输可能需要较长时间
   - 目录传输会递归传输所有子目录和文件
   - 确保有足够的磁盘空间

4. **跨平台兼容性：**
   - Windows 和 Unix 系统在路径分隔符上有差异，工具会自动处理
   - 终端大小检测在不同系统上可能略有差异

## 🐛 故障排除

### 连接失败

- 检查服务器地址和端口是否正确
- 确认网络连接正常
- 验证用户名和密码是否正确
- 检查服务器 SSH 服务是否运行

### 文件传输失败

- 确认本地和远程路径是否正确
- 检查文件/目录权限
- 确保目标目录存在或具有创建权限

### 命令执行失败

- 检查命令语法是否正确
- 确认用户在远程服务器上有执行权限
- 查看错误信息获取更多细节

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📧 联系方式

如有问题或建议，请提交 Issue。

---

**享受使用 GoSSH！** 🚀

