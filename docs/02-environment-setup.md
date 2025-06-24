# 第二章：环境准备

## 🛠️ Go 语言环境安装

### Windows 系统

1. **下载 Go 安装包**
   ```
   访问官网：https://golang.org/dl/
   下载：go1.24.x.windows-amd64.msi
   ```

2. **安装 Go**
   - 双击 msi 文件
   - 按照向导完成安装（默认安装到 C:\Go）
   - 安装程序会自动设置环境变量

3. **验证安装**
   ```cmd
   go version
   # 应该显示：go version go1.24.x windows/amd64
   ```

### macOS 系统

1. **使用 Homebrew 安装（推荐）**
   ```bash
   # 安装 Homebrew（如果还没有）
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   
   # 安装 Go
   brew install go
   ```

2. **或者下载安装包**
   ```
   访问官网：https://golang.org/dl/
   下载：go1.24.x.darwin-amd64.pkg
   双击安装
   ```

3. **验证安装**
   ```bash
   go version
   # 应该显示：go version go1.24.x darwin/amd64
   ```

### Linux 系统

1. **使用包管理器安装**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install golang-go
   
   # CentOS/RHEL
   sudo yum install golang
   
   # 或者使用 dnf
   sudo dnf install golang
   ```

2. **或者手动安装**
   ```bash
   # 下载
   wget https://golang.org/dl/go1.24.x.linux-amd64.tar.gz
   
   # 解压到 /usr/local
   sudo tar -C /usr/local -xzf go1.24.x.linux-amd64.tar.gz
   
   # 添加到 PATH
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **验证安装**
   ```bash
   go version
   # 应该显示：go version go1.24.x linux/amd64
   ```

## 📁 项目目录结构

### 创建项目目录

```bash
# 创建项目根目录
mkdir go-guess-game
cd go-guess-game

# 创建子目录
mkdir docs        # 文档目录
mkdir tests       # 测试相关文件（如果需要）
```

### 推荐的项目结构

```
go-guess-game/
├── main.go           # 主程序文件
├── main_test.go      # 测试文件
├── go.mod            # Go 模块文件
├── go.sum            # 依赖版本锁定文件（自动生成）
├── README.md         # 项目说明文档
├── docs/             # 详细文档目录
│   ├── README.md
│   ├── 01-project-overview.md
│   ├── 02-environment-setup.md
│   └── ...
└── guess-game        # 编译后的可执行文件
```

## 🚀 Go 模块初始化

### 什么是 Go 模块？

Go 模块（Go Modules）是 Go 1.11 引入的依赖管理系统，它解决了：
- 依赖版本管理
- 可重现的构建
- 项目隔离

### 初始化模块

```bash
# 在项目根目录执行
go mod init guess-number-game

# 这会创建 go.mod 文件
```

### go.mod 文件内容

```go
module guess-number-game

go 1.24
```

### 理解 go.mod 文件

- `module guess-number-game`：定义模块名称
- `go 1.24`：指定 Go 版本要求
- 随着项目发展，这里还会添加依赖项

## 🔧 开发工具配置

### VS Code 配置（推荐）

1. **安装 VS Code**
   ```
   访问：https://code.visualstudio.com/
   下载并安装适合您系统的版本
   ```

2. **安装 Go 扩展**
   - 打开 VS Code
   - 按 `Ctrl+Shift+X` 打开扩展面板
   - 搜索 "Go"
   - 安装 Google 官方的 Go 扩展

3. **配置 Go 工具**
   ```
   按 Ctrl+Shift+P 打开命令面板
   输入：Go: Install/Update Tools
   选择所有工具并安装
   ```

### 其他编辑器选择

1. **GoLand**：JetBrains 出品的专业 Go IDE
2. **Vim/Neovim**：配合 vim-go 插件
3. **Emacs**：配合 go-mode
4. **Sublime Text**：配合 GoSublime 插件

## 🧪 环境验证

### 创建测试文件

创建 `hello.go` 文件：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
    fmt.Println("环境配置成功！")
}
```

### 运行测试

```bash
# 直接运行
go run hello.go

# 编译后运行
go build hello.go
./hello        # Linux/macOS
hello.exe      # Windows
```

### 预期输出

```
Hello, Go!
环境配置成功！
```

## 📋 Go 工具链介绍

### 常用命令

```bash
# 运行程序
go run main.go

# 编译程序
go build main.go

# 编译并指定输出文件名
go build -o guess-game main.go

# 运行测试
go test

# 运行测试并显示详细信息
go test -v

# 运行基准测试
go test -bench=.

# 查看测试覆盖率
go test -cover

# 格式化代码
go fmt

# 检查代码问题
go vet

# 下载依赖
go mod download

# 清理未使用的依赖
go mod tidy
```

### Go 工具的作用

1. **go run**：直接运行 Go 源码，适合开发调试
2. **go build**：编译生成可执行文件，适合生产部署
3. **go test**：运行测试，确保代码质量
4. **go fmt**：自动格式化代码，保持代码风格一致
5. **go vet**：静态分析工具，发现潜在问题

## 🎯 环境配置检查清单

在继续下一章之前，请确保：

- [ ] Go 语言已正确安装（`go version` 命令有效）
- [ ] 项目目录已创建
- [ ] Go 模块已初始化（存在 go.mod 文件）
- [ ] 开发工具已配置（推荐 VS Code + Go 扩展）
- [ ] 能够成功运行简单的 Go 程序
- [ ] 了解基本的 Go 工具链命令

## 💡 常见问题解决

### 问题 1：go 命令找不到

**解决方案**：
```bash
# 检查 PATH 环境变量
echo $PATH

# 手动添加 Go 到 PATH（Linux/macOS）
export PATH=$PATH:/usr/local/go/bin

# 永久添加到 ~/.bashrc 或 ~/.zshrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

### 问题 2：GOPATH 相关错误

**解决方案**：
Go 1.11+ 使用模块系统，不再需要设置 GOPATH。如果遇到相关错误：
```bash
# 确保在项目目录内
pwd

# 确保有 go.mod 文件
ls go.mod

# 如果没有，重新初始化模块
go mod init your-project-name
```

### 问题 3：中文显示乱码

**解决方案**：
确保终端支持 UTF-8 编码：
- Windows：使用 PowerShell 或 Windows Terminal
- macOS/Linux：现代终端默认支持 UTF-8

## 🚀 下一步

环境准备完成后，在下一章中我们将：
- 设计游戏的核心数据结构
- 规划代码架构
- 定义接口和方法
- 开始编写第一个版本的代码

让我们继续前进！
