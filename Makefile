.PHONY: build run clean test fmt vet lint install watch help

# 构建可执行文件
build:
	@echo "正在构建..."
	@go build -o goss.exe .
	@echo "构建完成: goss.exe"

# 构建 Linux 版本
build-linux:
	@echo "正在构建 Linux 版本..."
	@GOOS=linux GOARCH=amd64 go build -o goss-linux .
	@echo "构建完成: goss-linux"

# 构建 macOS 版本
build-macos:
	@echo "正在构建 macOS 版本..."
	@GOOS=darwin GOARCH=amd64 go build -o goss-macos .
	@echo "构建完成: goss-macos"

# 直接运行程序（开发模式）
run:
	@go run main.go

# 运行并传递参数
run-interactive:
	@go run main.go interactive

run-connect:
	@go run main.go connect

run-list:
	@go run main.go list

# 清理构建文件
clean:
	@echo "正在清理..."
	@rm -f goss.exe goss-linux goss-macos
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@goimports -w .

# 代码检查
vet:
	@echo "运行 go vet..."
	@go vet ./...

# 代码检查（使用 golangci-lint）
lint:
	@echo "运行 linter..."
	@golangci-lint run ./...

# 安装依赖
install:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# 安装开发工具
install-tools:
	@echo "安装开发工具..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest

# 使用 air 热重载开发（需要先安装 air）
watch:
	@air

# 显示帮助
help:
	@echo "可用的命令:"
	@echo "  make build          - 构建可执行文件"
	@echo "  make build-linux    - 构建 Linux 版本"
	@echo "  make build-macos    - 构建 macOS 版本"
	@echo "  make run            - 运行程序"
	@echo "  make run-interactive - 运行交互式模式"
	@echo "  make clean          - 清理构建文件"
	@echo "  make test           - 运行测试"
	@echo "  make fmt            - 格式化代码"
	@echo "  make vet            - 代码检查"
	@echo "  make lint           - 运行 linter"
	@echo "  make install        - 安装依赖"
	@echo "  make install-tools  - 安装开发工具"
	@echo "  make watch          - 热重载开发模式（需要 air）"
