# 第四章：基础实现

## 🚀 开始编码

现在我们开始将设计转化为实际的代码。我们将按照以下顺序逐步实现：

1. 创建基本的项目结构
2. 实现 `Game` 结构体
3. 实现构造函数 `NewGame()`
4. 实现核心的游戏逻辑方法

## 📝 创建 main.go 文件

首先创建项目的主文件：

```go
package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)
```

### 导入包说明

- `bufio`：提供缓冲 I/O 功能，用于安全地读取用户输入
- `fmt`：格式化 I/O，用于打印和输入
- `math/rand`：随机数生成
- `os`：操作系统接口，用于访问标准输入
- `strconv`：字符串转换，用于将字符串转换为数字
- `strings`：字符串处理，用于去除空白字符等
- `time`：时间处理，用于设置随机数种子

## 🏗️ 定义 Game 结构体

```go
// Game 结构体用于管理游戏状态
type Game struct {
    targetNumber int           // 目标数字
    attempts     int           // 猜测次数
    scanner      *bufio.Scanner // 输入扫描器
}
```

### 结构体设计说明

1. **targetNumber**：存储当前游戏的目标数字（1-100）
2. **attempts**：记录玩家的猜测次数，用于最终统计
3. **scanner**：缓冲输入扫描器，提供安全的输入读取

## 🎲 实现构造函数

```go
// NewGame 创建新游戏实例
func NewGame() *Game {
    // 设置随机数种子
    rand.Seed(time.Now().UnixNano())
    
    return &Game{
        targetNumber: rand.Intn(100) + 1, // 生成 1-100 之间的随机数
        attempts:     0,
        scanner:      bufio.NewScanner(os.Stdin),
    }
}
```

### 实现要点

1. **随机数种子**：
   ```go
   rand.Seed(time.Now().UnixNano())
   ```
   - 使用当前时间的纳秒值作为种子
   - 确保每次运行程序都有不同的随机序列

2. **随机数生成**：
   ```go
   rand.Intn(100) + 1
   ```
   - `rand.Intn(100)` 生成 0-99 的随机数
   - `+1` 将范围调整为 1-100

3. **返回指针**：
   ```go
   return &Game{...}
   ```
   - 返回结构体指针而不是值
   - 避免大结构体的复制开销
   - 允许方法修改结构体内容

## 🎮 实现核心游戏逻辑

### checkGuess 方法

```go
// checkGuess 检查猜测结果
// 返回值：0 表示猜对，正数表示猜大了，负数表示猜小了
func (g *Game) checkGuess(guess int) int {
    if guess == g.targetNumber {
        return 0
    } else if guess > g.targetNumber {
        return 1
    } else {
        return -1
    }
}
```

### 方法设计说明

1. **方法接收者**：
   ```go
   func (g *Game) checkGuess(guess int) int
   ```
   - `(g *Game)` 表示这是 `Game` 结构体的方法
   - `g` 是接收者的名称，按惯例使用结构体名的首字母
   - 使用指针接收者，虽然这个方法不修改状态，但保持一致性

2. **返回值设计**：
   - `0`：猜对了
   - `1`：猜大了（guess > target）
   - `-1`：猜小了（guess < target）
   - 这种设计简洁明了，易于理解和测试

## 📥 实现输入处理

```go
// getPlayerGuess 获取玩家的猜测输入
func (g *Game) getPlayerGuess() (int, error) {
    fmt.Print("请输入你的猜测：")
    
    if !g.scanner.Scan() {
        return 0, fmt.Errorf("读取输入失败")
    }
    
    input := strings.TrimSpace(g.scanner.Text())
    if input == "" {
        return 0, fmt.Errorf("输入不能为空")
    }
    
    guess, err := strconv.Atoi(input)
    if err != nil {
        return 0, fmt.Errorf("请输入一个有效的数字")
    }
    
    if guess < 1 || guess > 100 {
        return 0, fmt.Errorf("数字必须在 1-100 之间")
    }
    
    return guess, nil
}
```

### 输入处理流程

1. **显示提示**：
   ```go
   fmt.Print("请输入你的猜测：")
   ```

2. **读取输入**：
   ```go
   if !g.scanner.Scan() {
       return 0, fmt.Errorf("读取输入失败")
   }
   ```
   - `scanner.Scan()` 读取一行输入
   - 返回 `false` 表示读取失败（如 EOF 或错误）

3. **处理空白**：
   ```go
   input := strings.TrimSpace(g.scanner.Text())
   ```
   - 去除输入前后的空白字符
   - 防止用户意外输入空格导致错误

4. **验证非空**：
   ```go
   if input == "" {
       return 0, fmt.Errorf("输入不能为空")
   }
   ```

5. **转换为数字**：
   ```go
   guess, err := strconv.Atoi(input)
   if err != nil {
       return 0, fmt.Errorf("请输入一个有效的数字")
   }
   ```
   - `strconv.Atoi()` 将字符串转换为整数
   - 如果转换失败，返回错误

6. **范围验证**：
   ```go
   if guess < 1 || guess > 100 {
       return 0, fmt.Errorf("数字必须在 1-100 之间")
   }
   ```

### 错误处理策略

我们使用 Go 语言的标准错误处理模式：
- 函数返回 `(结果, error)` 的形式
- 成功时返回 `(有效值, nil)`
- 失败时返回 `(零值, 错误信息)`

## 🎯 实现游戏主循环

```go
// Start 开始游戏
func (g *Game) Start() {
    fmt.Println("欢迎来到猜数字游戏！")
    fmt.Println("我已经想好了一个 1-100 之间的数字，请开始猜测：")
    
    for {
        guess, err := g.getPlayerGuess()
        if err != nil {
            fmt.Printf("输入错误：%v，请重新输入。\n", err)
            continue
        }
        
        g.attempts++
        result := g.checkGuess(guess)
        
        if result == 0 {
            // 猜对了
            fmt.Printf("恭喜你！猜对了！你总共猜了 %d 次。\n", g.attempts)
            break
        } else if result > 0 {
            fmt.Println("太大了！请再试一次。")
        } else {
            fmt.Println("太小了！请再试一次。")
        }
    }
}
```

### 主循环逻辑

1. **显示欢迎信息**：
   ```go
   fmt.Println("欢迎来到猜数字游戏！")
   fmt.Println("我已经想好了一个 1-100 之间的数字，请开始猜测：")
   ```

2. **无限循环**：
   ```go
   for {
       // 游戏逻辑
   }
   ```
   - 使用无限循环，直到猜对为止
   - 通过 `break` 退出循环

3. **获取输入并处理错误**：
   ```go
   guess, err := g.getPlayerGuess()
   if err != nil {
       fmt.Printf("输入错误：%v，请重新输入。\n", err)
       continue
   }
   ```
   - 如果输入有错误，显示错误信息并继续循环
   - 不会因为输入错误而退出游戏

4. **更新统计**：
   ```go
   g.attempts++
   ```
   - 只有在输入有效时才增加猜测次数

5. **检查结果并给出反馈**：
   ```go
   result := g.checkGuess(guess)
   
   if result == 0 {
       // 猜对了，显示成功信息并退出
   } else if result > 0 {
       // 猜大了
   } else {
       // 猜小了
   }
   ```

## 🔧 测试基础功能

在继续实现完整程序之前，我们可以创建一个简单的测试来验证基础功能：

```go
// 临时的 main 函数用于测试
func main() {
    game := NewGame()
    fmt.Printf("目标数字：%d\n", game.targetNumber) // 调试用，实际游戏中不显示
    
    // 测试 checkGuess 方法
    fmt.Printf("猜测 50，结果：%d\n", game.checkGuess(50))
    fmt.Printf("猜测 %d，结果：%d\n", game.targetNumber, game.checkGuess(game.targetNumber))
}
```

运行测试：
```bash
go run main.go
```

## 📊 代码质量检查

### 使用 Go 工具

1. **格式化代码**：
   ```bash
   go fmt main.go
   ```

2. **检查代码问题**：
   ```bash
   go vet main.go
   ```

3. **运行程序**：
   ```bash
   go run main.go
   ```

### 代码审查要点

1. **命名规范**：
   - 结构体和方法使用驼峰命名
   - 私有成员小写开头，公有成员大写开头

2. **注释规范**：
   - 每个公有类型和方法都有注释
   - 注释以类型或方法名开头

3. **错误处理**：
   - 所有可能的错误都被处理
   - 错误信息清晰明确

## 🎯 本章总结

在这一章中，我们实现了：

1. ✅ **项目结构**：创建了基本的 Go 项目文件
2. ✅ **数据结构**：定义了 `Game` 结构体
3. ✅ **构造函数**：实现了 `NewGame()` 方法
4. ✅ **核心逻辑**：实现了 `checkGuess()` 方法
5. ✅ **输入处理**：实现了 `getPlayerGuess()` 方法
6. ✅ **游戏循环**：实现了 `Start()` 方法

### 当前代码特点

- **结构清晰**：每个方法职责单一
- **错误处理完善**：不会因为错误输入而崩溃
- **用户友好**：提供清晰的中文提示
- **易于测试**：方法独立，便于单元测试

## 🚀 下一步

在下一章中，我们将：
- 完善用户交互界面
- 添加游戏继续功能
- 实现完整的 main 函数
- 优化用户体验

我们的基础框架已经搭建完成，继续前进！
