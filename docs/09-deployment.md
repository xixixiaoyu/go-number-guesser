# 第九章：部署运行

## 🚀 编译与构建

Go 语言的一个重要优势是能够编译成单一的可执行文件，无需依赖外部运行时环境。在这一章中，我们将学习如何编译、打包和分发我们的猜数字游戏。

## 🔨 基本编译

### 简单编译

```bash
# 编译当前目录的 Go 程序
go build

# 编译并指定输出文件名
go build -o guess-game

# 编译指定文件
go build main.go

# 编译并指定输出文件名和路径
go build -o bin/guess-game main.go
```

### 编译选项详解

```bash
# 显示编译过程
go build -v

# 编译时显示更多信息
go build -x

# 禁用优化（调试用）
go build -gcflags="-N -l"

# 减小可执行文件大小
go build -ldflags="-s -w"
```

#### 编译选项说明

1. **`-v`**：显示被编译的包名
2. **`-x`**：显示执行的命令
3. **`-gcflags="-N -l"`**：
   - `-N`：禁用优化
   - `-l`：禁用内联
   - 用于调试，保留更多调试信息

4. **`-ldflags="-s -w"`**：
   - `-s`：去掉符号表
   - `-w`：去掉调试信息
   - 可以显著减小可执行文件大小

## 🌍 跨平台编译

Go 语言支持交叉编译，可以在一个平台上编译出其他平台的可执行文件。

### 查看支持的平台

```bash
# 查看当前环境
go env GOOS GOARCH

# 查看所有支持的平台
go tool dist list
```

### 跨平台编译示例

```bash
# 编译 Windows 64位版本
GOOS=windows GOARCH=amd64 go build -o guess-game.exe

# 编译 Linux 64位版本
GOOS=linux GOARCH=amd64 go build -o guess-game-linux

# 编译 macOS 64位版本
GOOS=darwin GOARCH=amd64 go build -o guess-game-macos

# 编译 ARM64 版本（如 Apple M1）
GOOS=darwin GOARCH=arm64 go build -o guess-game-macos-arm64

# 编译 Linux ARM 版本（如树莓派）
GOOS=linux GOARCH=arm go build -o guess-game-linux-arm
```

### 批量编译脚本

创建 `build.sh` 脚本：

```bash
#!/bin/bash

# 项目名称
PROJECT_NAME="guess-game"

# 版本号
VERSION="v1.0.0"

# 创建输出目录
mkdir -p dist

# 编译不同平台版本
echo "开始编译..."

# Windows 64位
echo "编译 Windows 64位版本..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/${PROJECT_NAME}-${VERSION}-windows-amd64.exe

# Linux 64位
echo "编译 Linux 64位版本..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/${PROJECT_NAME}-${VERSION}-linux-amd64

# macOS 64位
echo "编译 macOS 64位版本..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/${PROJECT_NAME}-${VERSION}-darwin-amd64

# macOS ARM64 (Apple Silicon)
echo "编译 macOS ARM64版本..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/${PROJECT_NAME}-${VERSION}-darwin-arm64

echo "编译完成！输出目录：dist/"
```

运行脚本：

```bash
chmod +x build.sh
./build.sh
```

## 📦 打包与分发

### 创建发布包

```bash
# 创建发布目录结构
mkdir -p release/guess-game-v1.0.0
cd release/guess-game-v1.0.0

# 复制必要文件
cp ../../dist/guess-game-v1.0.0-* .
cp ../../README.md .
cp ../../LICENSE .

# 创建安装说明
cat > INSTALL.md << 'EOF'
# 安装说明

## Windows 用户
1. 下载 `guess-game-v1.0.0-windows-amd64.exe`
2. 双击运行或在命令行中执行

## Linux 用户
1. 下载 `guess-game-v1.0.0-linux-amd64`
2. 添加执行权限：`chmod +x guess-game-v1.0.0-linux-amd64`
3. 运行：`./guess-game-v1.0.0-linux-amd64`

## macOS 用户
1. 下载对应版本：
   - Intel Mac: `guess-game-v1.0.0-darwin-amd64`
   - Apple Silicon: `guess-game-v1.0.0-darwin-arm64`
2. 添加执行权限：`chmod +x guess-game-v1.0.0-darwin-*`
3. 运行程序
EOF

# 创建压缩包
cd ..
tar -czf guess-game-v1.0.0.tar.gz guess-game-v1.0.0/
zip -r guess-game-v1.0.0.zip guess-game-v1.0.0/
```

### Docker 容器化

创建 `Dockerfile`：

```dockerfile
# 使用官方 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o guess-game

# 使用轻量级镜像作为运行环境
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /home/appuser

# 从构建阶段复制可执行文件
COPY --from=builder /app/guess-game .

# 更改文件所有者
RUN chown appuser:appuser guess-game

# 切换到非 root 用户
USER appuser

# 运行应用
CMD ["./guess-game"]
```

构建和运行 Docker 镜像：

```bash
# 构建镜像
docker build -t guess-game:v1.0.0 .

# 运行容器
docker run -it guess-game:v1.0.0
```

## 🎯 版本管理

### 语义化版本

遵循语义化版本规范（Semantic Versioning）：

```
主版本号.次版本号.修订号

例如：
1.0.0 - 初始版本
1.0.1 - 修复 bug
1.1.0 - 添加新功能
2.0.0 - 重大更新，可能不兼容
```

### 在代码中嵌入版本信息

```go
package main

import "fmt"

var (
    Version   = "dev"      // 版本号，构建时注入
    BuildTime = "unknown"  // 构建时间，构建时注入
    GitCommit = "unknown"  // Git 提交哈希，构建时注入
)

func showVersion() {
    fmt.Printf("猜数字游戏 %s\n", Version)
    fmt.Printf("构建时间: %s\n", BuildTime)
    fmt.Printf("Git 提交: %s\n", GitCommit)
}
```

### 构建时注入版本信息

```bash
# 获取版本信息
VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT=$(git rev-parse HEAD)

# 编译时注入版本信息
go build -ldflags="-X main.Version=${VERSION} -X 'main.BuildTime=${BUILD_TIME}' -X main.GitCommit=${GIT_COMMIT}" -o guess-game
```

### Makefile 自动化

创建 `Makefile`：

```makefile
# 项目配置
PROJECT_NAME := guess-game
VERSION := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse HEAD)

# 构建标志
LDFLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X 'main.BuildTime=$(BUILD_TIME)' -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
.PHONY: all
all: clean test build

# 清理
.PHONY: clean
clean:
	rm -rf dist/
	rm -rf release/

# 测试
.PHONY: test
test:
	go test -v -cover

# 本地构建
.PHONY: build
build:
	go build $(LDFLAGS) -o $(PROJECT_NAME)

# 跨平台构建
.PHONY: build-all
build-all: clean
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(PROJECT_NAME)-$(VERSION)-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(PROJECT_NAME)-$(VERSION)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(PROJECT_NAME)-$(VERSION)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(PROJECT_NAME)-$(VERSION)-darwin-arm64

# 创建发布包
.PHONY: release
release: build-all
	mkdir -p release
	cd dist && for file in *; do \
		mkdir -p ../release/$$file && \
		cp $$file ../release/$$file/ && \
		cp ../README.md ../release/$$file/ && \
		cp ../LICENSE ../release/$$file/ && \
		cd ../release && tar -czf $$file.tar.gz $$file/ && cd ../dist; \
	done

# 安装到本地
.PHONY: install
install: build
	cp $(PROJECT_NAME) /usr/local/bin/

# 运行
.PHONY: run
run:
	go run main.go

# 格式化代码
.PHONY: fmt
fmt:
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	go vet ./...

# 显示帮助
.PHONY: help
help:
	@echo "可用的命令："
	@echo "  make build      - 本地构建"
	@echo "  make build-all  - 跨平台构建"
	@echo "  make test       - 运行测试"
	@echo "  make clean      - 清理文件"
	@echo "  make release    - 创建发布包"
	@echo "  make install    - 安装到系统"
	@echo "  make run        - 运行程序"
	@echo "  make fmt        - 格式化代码"
	@echo "  make vet        - 代码检查"
```

使用 Makefile：

```bash
# 查看帮助
make help

# 运行测试和构建
make

# 跨平台构建
make build-all

# 创建发布包
make release
```

## 📋 部署检查清单

### 构建前检查

- [ ] 代码已提交到版本控制系统
- [ ] 所有测试都通过
- [ ] 代码已经过格式化和静态检查
- [ ] 版本号已更新
- [ ] 文档已更新

### 构建检查

- [ ] 本地构建成功
- [ ] 跨平台构建成功
- [ ] 可执行文件大小合理
- [ ] 版本信息正确嵌入

### 发布前检查

- [ ] 在不同平台上测试运行
- [ ] 检查文件权限
- [ ] 验证压缩包完整性
- [ ] 确认发布说明准确

## 🔧 常见问题解决

### 问题 1：可执行文件过大

**解决方案**：
```bash
# 使用 ldflags 减小文件大小
go build -ldflags="-s -w" -o guess-game

# 使用 UPX 进一步压缩（可选）
upx --best guess-game
```

### 问题 2：跨平台编译失败

**解决方案**：
```bash
# 确保目标平台支持
go tool dist list | grep linux

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

### 问题 3：Docker 构建慢

**解决方案**：
```dockerfile
# 使用多阶段构建
# 利用 Docker 缓存层
# 先复制 go.mod，再复制源码
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

## 🎯 本章总结

在这一章中，我们学习了项目的部署和分发：

1. ✅ **基本编译**：掌握了 Go 编译的基本用法
2. ✅ **跨平台编译**：学会了为不同平台编译程序
3. ✅ **打包分发**：了解了如何创建发布包
4. ✅ **容器化**：学习了 Docker 容器化部署
5. ✅ **版本管理**：掌握了版本信息的管理方法
6. ✅ **自动化构建**：使用 Makefile 自动化构建流程

### 关键收获

- **单一可执行文件**：Go 编译的优势
- **跨平台支持**：一次编写，到处运行
- **自动化构建**：提高开发效率
- **版本管理**：规范的版本控制

## 🚀 下一步

在最后一章中，我们将：
- 探讨功能扩展的可能性
- 学习性能优化技巧
- 了解 Go 语言的高级特性
- 规划后续学习路径

部署运行章节完成，让我们继续探索扩展功能！
