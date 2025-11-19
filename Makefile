.PHONY: build run clean test help

# 变量
BINARY_NAME=webhook-listener
MAIN_PATH=./cmd/webhook-listener
BUILD_DIR=./bin
CONFIG_FILE=config.yaml
PORT=8375

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  make build     - 构建可执行文件"
	@echo "  make run       - 运行程序（需要先创建 config.yaml）"
	@echo "  make clean     - 清理构建文件"
	@echo "  make test      - 运行测试"
	@echo "  make install   - 安装到系统路径"

# 构建
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行
run:
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "错误: 配置文件 $(CONFIG_FILE) 不存在"; \
		echo "请先复制 config/config.example.yaml 到 $(CONFIG_FILE) 并修改配置"; \
		exit 1; \
	fi
	@go run $(MAIN_PATH) -config $(CONFIG_FILE) -port $(PORT)

# 清理
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "清理完成"

# 测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 安装
install: build
	@echo "安装到 /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成"

# 交叉编译（Linux）
build-linux:
	@echo "构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

# 交叉编译（macOS）
build-darwin:
	@echo "构建 macOS 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "构建完成"

# 交叉编译（Windows）
build-windows:
	@echo "构建 Windows 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# 构建所有平台
build-all: build-linux build-darwin build-windows
	@echo "所有平台构建完成"

