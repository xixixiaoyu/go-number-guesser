# 第五章：用户交互

## 🎨 用户界面设计

一个好的命令行程序不仅要功能正确，还要有良好的用户体验。在这一章中，我们将完善用户交互界面，让游戏更加友好和美观。

## 🎯 设计目标

1. **视觉美观**：使用分隔线、表情符号等元素
2. **信息清晰**：提示信息明确，反馈及时
3. **操作简单**：输入要求简单，容错性强
4. **体验流畅**：界面布局合理，流程自然

## 🖼️ 界面布局设计

### 游戏标题界面

我们在 main 函数中直接实现标题显示：

```go
fmt.Println(strings.Repeat("=", 50))
fmt.Println("🎯 Go 语言猜数字游戏")
fmt.Println(strings.Repeat("=", 50))
```

**设计说明**：
- 使用 50 个等号创建醒目的分隔线
- 添加表情符号增加趣味性
- 标题居中显示，视觉效果好

### 新游戏开始界面

```go
fmt.Println("\n" + strings.Repeat("=", 30))
fmt.Println("开始新游戏！")
fmt.Println(strings.Repeat("=", 30))
```

**设计说明**：
- 使用较短的分隔线区分不同游戏轮次
- 添加空行增加视觉间隔
- 简洁明了的提示信息

## 🔄 实现游戏继续功能

### askContinue 函数

```go
// askContinue 询问是否继续游戏
func askContinue() bool {
    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("是否继续游戏？(y/n)：")
        
        if !scanner.Scan() {
            fmt.Println("读取输入失败，默认退出游戏。")
            return false
        }
        
        input := strings.ToLower(strings.TrimSpace(scanner.Text()))
        
        switch input {
        case "y", "yes", "是":
            return true
        case "n", "no", "否":
            return false
        default:
            fmt.Println("请输入 y(是) 或 n(否)。")
        }
    }
}
```

### 实现要点分析

1. **多种输入支持**：
   ```go
   case "y", "yes", "是":
       return true
   case "n", "no", "否":
       return false
   ```
   - 支持英文和中文输入
   - 支持完整单词和简写
   - 提高用户体验的灵活性

2. **输入预处理**：
   ```go
   input := strings.ToLower(strings.TrimSpace(scanner.Text()))
   ```
   - `TrimSpace()` 去除前后空白
   - `ToLower()` 转换为小写，忽略大小写

3. **错误处理**：
   ```go
   if !scanner.Scan() {
       fmt.Println("读取输入失败，默认退出游戏。")
       return false
   }
   ```
   - 处理输入读取失败的情况
   - 提供默认行为，避免程序卡死

4. **输入验证循环**：
   ```go
   for {
       // 获取输入
       // 验证输入
       // 如果无效，提示并继续循环
   }
   ```
   - 持续询问直到获得有效输入
   - 提供清晰的错误提示

## 🎮 完整的 main 函数

```go
// main 主函数
func main() {
    fmt.Println(strings.Repeat("=", 50))
    fmt.Println("🎯 Go 语言猜数字游戏")
    fmt.Println(strings.Repeat("=", 50))

    for {
        game := NewGame()
        game.Start()

        if !askContinue() {
            fmt.Println("感谢游戏！再见！👋")
            break
        }

        fmt.Println("\n" + strings.Repeat("=", 30))
        fmt.Println("开始新游戏！")
        fmt.Println(strings.Repeat("=", 30))
    }
}
```

### 主函数流程分析

1. **显示游戏标题**：
   ```go
   fmt.Println(strings.Repeat("=", 50))
   fmt.Println("🎯 Go 语言猜数字游戏")
   fmt.Println(strings.Repeat("=", 50))
   ```

2. **游戏主循环**：
   ```go
   for {
       // 创建新游戏
       // 开始游戏
       // 询问是否继续
       // 如果不继续，退出循环
   }
   ```

3. **游戏实例管理**：
   ```go
   game := NewGame()
   game.Start()
   ```
   - 每轮游戏创建新的实例
   - 确保游戏状态独立

4. **优雅退出**：
   ```go
   if !askContinue() {
       fmt.Println("感谢游戏！再见！👋")
       break
   }
   ```
   - 友好的结束信息
   - 使用表情符号增加亲和力

## 💬 改进反馈信息

### 优化 Start 方法的反馈

```go
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

### 反馈信息设计原则

1. **积极正面**：
   - "恭喜你！猜对了！" 而不是 "正确"
   - 使用感叹号增加情感色彩

2. **信息完整**：
   - "你总共猜了 X 次" 提供统计信息
   - 让玩家了解自己的表现

3. **指导性强**：
   - "太大了！请再试一次。" 明确指出调整方向
   - "请重新输入" 告诉用户下一步操作

## 🎨 界面美化技巧

### 使用分隔线

```go
// 主标题分隔线
strings.Repeat("=", 50)

// 次级标题分隔线
strings.Repeat("=", 30)

// 简单分隔线
strings.Repeat("-", 20)
```

### 合理使用空行

```go
// 在重要信息前后添加空行
fmt.Println("\n" + strings.Repeat("=", 30))
fmt.Println("开始新游戏！")
fmt.Println(strings.Repeat("=", 30))
```

### 表情符号的使用

```go
// 游戏标题
"🎯 Go 语言猜数字游戏"

// 结束信息
"感谢游戏！再见！👋"

// 可以考虑的其他表情符号
"🎲" // 骰子，表示随机
"🎉" // 庆祝，表示成功
"❌" // 错误标记
"✅" // 正确标记
```

## 🔧 输入体验优化

### 改进输入提示

```go
func (g *Game) getPlayerGuess() (int, error) {
    fmt.Print("请输入你的猜测：")
    // ... 其他代码保持不变
}
```

**优化建议**：
- 使用 `fmt.Print()` 而不是 `fmt.Println()`，让输入在同一行
- 提示信息简洁明了
- 保持一致的提示格式

### 错误信息优化

```go
// 具体的错误信息
return 0, fmt.Errorf("输入不能为空")
return 0, fmt.Errorf("请输入一个有效的数字")
return 0, fmt.Errorf("数字必须在 1-100 之间")
```

**优化特点**：
- 错误信息具体明确
- 使用中文，用户友好
- 告诉用户正确的输入格式

## 🧪 交互测试

### 测试用例设计

1. **正常流程测试**：
   - 输入有效数字
   - 多次猜测直到成功
   - 选择继续或退出

2. **异常输入测试**：
   - 输入非数字字符
   - 输入超出范围的数字
   - 输入空值
   - 输入空格

3. **边界条件测试**：
   - 输入 1 和 100
   - 一次猜中
   - 多次猜测

### 手动测试步骤

```bash
# 编译并运行
go build -o guess-game main.go
./guess-game

# 测试各种输入情况
# 1. 正常数字：50
# 2. 非数字：abc
# 3. 超出范围：150
# 4. 空输入：直接回车
# 5. 继续游戏：y/n/yes/no/是/否
```

## 📊 用户体验评估

### 评估标准

1. **易用性**：
   - [ ] 提示信息清晰
   - [ ] 输入要求简单
   - [ ] 错误处理友好

2. **美观性**：
   - [ ] 界面布局合理
   - [ ] 视觉元素协调
   - [ ] 信息层次分明

3. **功能性**：
   - [ ] 所有功能正常工作
   - [ ] 错误不会导致崩溃
   - [ ] 流程逻辑正确

## 🎯 本章总结

在这一章中，我们完善了用户交互功能：

1. ✅ **界面美化**：添加了分隔线和表情符号
2. ✅ **游戏继续功能**：实现了 `askContinue()` 函数
3. ✅ **完整主函数**：实现了多轮游戏循环
4. ✅ **反馈优化**：改进了提示和错误信息
5. ✅ **输入体验**：支持多种输入格式

### 当前功能特点

- **用户友好**：中文界面，清晰提示
- **容错性强**：完善的错误处理
- **视觉美观**：合理的界面布局
- **操作简单**：直观的交互流程

## 🚀 下一步

在下一章中，我们将：
- 深入分析游戏逻辑的实现细节
- 优化算法性能
- 添加更多的游戏特性
- 完善状态管理

用户交互部分已经完成，让我们继续完善游戏逻辑！
