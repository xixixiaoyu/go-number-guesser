# 第七章：测试编写

## 🧪 测试的重要性

测试是软件开发中不可或缺的一环，特别是在 Go 语言中，测试被视为一等公民。好的测试不仅能发现 bug，还能：

1. **文档化代码行为**：测试用例就是最好的使用示例
2. **支持重构**：有了测试保护，可以安全地修改代码
3. **提高代码质量**：编写可测试的代码通常质量更高
4. **增强信心**：完善的测试让开发者更有信心发布代码

## 📋 Go 测试基础

### 测试文件命名规则

Go 语言的测试文件必须遵循特定的命名规则：

```
main.go      -> main_test.go
game.go      -> game_test.go
utils.go     -> utils_test.go
```

**规则**：
- 测试文件名必须以 `_test.go` 结尾
- 测试文件与被测试文件在同一个包中

### 测试函数命名规则

```go
func TestFunctionName(t *testing.T) { ... }    // 单元测试
func BenchmarkFunctionName(b *testing.B) { ... } // 基准测试
func ExampleFunctionName() { ... }              // 示例测试
```

## 🔬 单元测试实现

### 测试 NewGame 函数

```go
package main

import (
    "testing"
)

// TestNewGame 测试游戏创建
func TestNewGame(t *testing.T) {
    game := NewGame()
    
    if game == nil {
        t.Fatal("NewGame() 返回了 nil")
    }
    
    if game.targetNumber < 1 || game.targetNumber > 100 {
        t.Errorf("目标数字 %d 不在 1-100 范围内", game.targetNumber)
    }
    
    if game.attempts != 0 {
        t.Errorf("初始猜测次数应该为 0，实际为 %d", game.attempts)
    }
}
```

#### 测试要点分析

1. **空值检查**：
   ```go
   if game == nil {
       t.Fatal("NewGame() 返回了 nil")
   }
   ```
   - 使用 `t.Fatal()` 而不是 `t.Error()`
   - 因为如果 game 为 nil，后续测试会 panic

2. **范围验证**：
   ```go
   if game.targetNumber < 1 || game.targetNumber > 100 {
       t.Errorf("目标数字 %d 不在 1-100 范围内", game.targetNumber)
   }
   ```
   - 验证随机数生成的正确性
   - 使用 `t.Errorf()` 提供详细的错误信息

3. **初始状态检查**：
   ```go
   if game.attempts != 0 {
       t.Errorf("初始猜测次数应该为 0，实际为 %d", game.attempts)
   }
   ```

### 测试 checkGuess 方法

```go
// TestCheckGuess 测试猜测检查逻辑
func TestCheckGuess(t *testing.T) {
    game := &Game{targetNumber: 50, attempts: 0}
    
    // 测试猜对
    result := game.checkGuess(50)
    if result != 0 {
        t.Errorf("猜对时应该返回 0，实际返回 %d", result)
    }
    
    // 测试猜大了
    result = game.checkGuess(75)
    if result <= 0 {
        t.Errorf("猜大时应该返回正数，实际返回 %d", result)
    }
    
    // 测试猜小了
    result = game.checkGuess(25)
    if result >= 0 {
        t.Errorf("猜小时应该返回负数，实际返回 %d", result)
    }
}
```

#### 测试设计思路

1. **创建测试实例**：
   ```go
   game := &Game{targetNumber: 50, attempts: 0}
   ```
   - 直接创建结构体实例，不依赖 `NewGame()`
   - 使用固定的目标数字，确保测试结果可预测

2. **覆盖所有分支**：
   - 猜对：`guess == target`
   - 猜大：`guess > target`
   - 猜小：`guess < target`

3. **验证返回值**：
   - 不仅检查返回值的正确性
   - 还检查返回值的符号（正数、负数、零）

### 测试随机性

```go
// TestRandomness 测试随机数生成的多样性
func TestRandomness(t *testing.T) {
    numbers := make(map[int]bool)
    
    // 生成多个游戏实例，检查目标数字的多样性
    for i := 0; i < 20; i++ {
        game := NewGame()
        numbers[game.targetNumber] = true
    }
    
    // 至少应该有几个不同的数字（考虑到随机性，不要求太严格）
    if len(numbers) < 3 {
        t.Errorf("随机数生成缺乏多样性，20次生成只有 %d 个不同数字", len(numbers))
    }
}
```

#### 随机性测试的挑战

1. **不确定性**：随机测试结果不确定
2. **平衡性**：既要检测随机性，又不能过于严格
3. **重复性**：测试应该能够重复运行

#### 解决方案

1. **统计方法**：使用大样本检测分布
2. **合理阈值**：设置合理的期望值
3. **多次运行**：在 CI/CD 中多次运行测试

## 📊 基准测试

### 基准测试的作用

1. **性能测量**：测量函数的执行时间
2. **性能对比**：比较不同实现的性能
3. **性能回归检测**：确保性能不会退化

### 实现基准测试

```go
// BenchmarkNewGame 基准测试游戏创建性能
func BenchmarkNewGame(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewGame()
    }
}

// BenchmarkCheckGuess 基准测试猜测检查性能
func BenchmarkCheckGuess(b *testing.B) {
    game := &Game{targetNumber: 50}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        game.checkGuess(25)
    }
}
```

#### 基准测试要点

1. **循环结构**：
   ```go
   for i := 0; i < b.N; i++ {
       // 被测试的代码
   }
   ```
   - `b.N` 由测试框架自动调整
   - 框架会运行足够多次以获得稳定的测量结果

2. **重置计时器**：
   ```go
   b.ResetTimer()
   ```
   - 在准备工作完成后重置计时器
   - 确保只测量核心逻辑的性能

3. **避免编译器优化**：
   ```go
   // 如果需要使用结果，可以赋值给包级变量
   var result int
   
   func BenchmarkCheckGuess(b *testing.B) {
       game := &Game{targetNumber: 50}
       var r int
       
       b.ResetTimer()
       for i := 0; i < b.N; i++ {
           r = game.checkGuess(25)
       }
       result = r // 防止编译器优化掉函数调用
   }
   ```

## 🎯 表格驱动测试

### 什么是表格驱动测试？

表格驱动测试是 Go 语言中常用的测试模式，通过定义测试用例表格来批量执行测试。

### 实现表格驱动测试

```go
func TestCheckGuessTableDriven(t *testing.T) {
    game := &Game{targetNumber: 50}
    
    tests := []struct {
        name     string
        guess    int
        expected int
    }{
        {"猜对了", 50, 0},
        {"猜大了", 75, 1},
        {"猜小了", 25, -1},
        {"边界值-最小", 1, -1},
        {"边界值-最大", 100, 1},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := game.checkGuess(tt.guess)
            if result != tt.expected {
                t.Errorf("checkGuess(%d) = %d, 期望 %d", 
                    tt.guess, result, tt.expected)
            }
        })
    }
}
```

#### 表格驱动测试的优势

1. **清晰的测试用例**：所有测试用例一目了然
2. **易于添加用例**：只需在表格中添加新行
3. **子测试支持**：每个用例都是独立的子测试
4. **详细的错误报告**：可以精确定位失败的用例

## 🔍 测试覆盖率

### 什么是测试覆盖率？

测试覆盖率是衡量测试质量的重要指标，表示测试执行了多少代码。

### 查看测试覆盖率

```bash
# 运行测试并生成覆盖率报告
go test -cover

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out

# 查看 HTML 格式的覆盖率报告
go tool cover -html=coverage.out
```

### 覆盖率分析

```bash
# 示例输出
PASS
coverage: 85.7% of statements
ok      guess-number-game    0.002s
```

#### 覆盖率指标

1. **语句覆盖率**：执行了多少语句
2. **分支覆盖率**：执行了多少分支
3. **函数覆盖率**：调用了多少函数

#### 提高覆盖率的方法

1. **测试所有分支**：确保 if-else 的每个分支都被测试
2. **测试边界条件**：测试边界值和异常情况
3. **测试错误路径**：不仅测试正常流程，还要测试错误处理

## 🧪 测试最佳实践

### 1. 测试命名

```go
// 好的命名
func TestNewGame(t *testing.T) { ... }
func TestCheckGuess_WhenGuessIsCorrect(t *testing.T) { ... }
func TestCheckGuess_WhenGuessIsTooHigh(t *testing.T) { ... }

// 不好的命名
func Test1(t *testing.T) { ... }
func TestStuff(t *testing.T) { ... }
```

### 2. 测试结构

```go
func TestFunction(t *testing.T) {
    // Arrange - 准备测试数据
    game := &Game{targetNumber: 50}
    
    // Act - 执行被测试的操作
    result := game.checkGuess(75)
    
    // Assert - 验证结果
    if result <= 0 {
        t.Errorf("期望正数，得到 %d", result)
    }
}
```

### 3. 错误信息

```go
// 好的错误信息
t.Errorf("checkGuess(%d) = %d, 期望 %d", guess, result, expected)

// 不好的错误信息
t.Error("测试失败")
```

### 4. 测试独立性

```go
// 每个测试都应该独立
func TestA(t *testing.T) {
    // 不依赖其他测试的结果
}

func TestB(t *testing.T) {
    // 不依赖 TestA 的执行
}
```

## 🚀 运行测试

### 基本命令

```bash
# 运行所有测试
go test

# 运行测试并显示详细信息
go test -v

# 运行特定的测试
go test -run TestNewGame

# 运行基准测试
go test -bench=.

# 运行基准测试并显示内存分配
go test -bench=. -benchmem

# 查看测试覆盖率
go test -cover
```

### 高级选项

```bash
# 并行运行测试
go test -parallel 4

# 设置超时时间
go test -timeout 30s

# 运行测试多次
go test -count 10

# 竞态检测
go test -race
```

## 📊 测试报告解读

### 单元测试报告

```
=== RUN   TestNewGame
--- PASS: TestNewGame (0.00s)
=== RUN   TestCheckGuess
--- PASS: TestCheckGuess (0.00s)
=== RUN   TestRandomness
--- PASS: TestRandomness (0.00s)
PASS
coverage: 85.7% of statements
ok      guess-number-game    0.002s
```

### 基准测试报告

```
BenchmarkNewGame-8         1000000      1023 ns/op
BenchmarkCheckGuess-8     50000000        32.1 ns/op
```

**解读**：
- `BenchmarkNewGame-8`：函数名和 GOMAXPROCS
- `1000000`：执行次数
- `1023 ns/op`：每次操作的平均时间

## 🎯 本章总结

在这一章中，我们学习了 Go 语言的测试技术：

1. ✅ **单元测试**：编写了完整的测试用例
2. ✅ **基准测试**：测量了函数的性能
3. ✅ **表格驱动测试**：学习了高效的测试模式
4. ✅ **测试覆盖率**：了解了如何衡量测试质量
5. ✅ **测试最佳实践**：掌握了测试编写的规范

### 关键收获

- **测试是文档**：好的测试就是最好的使用说明
- **测试驱动开发**：先写测试，再写实现
- **持续测试**：每次修改代码都要运行测试
- **质量保证**：测试是代码质量的重要保障

## 🚀 下一步

在下一章中，我们将：
- 完善项目文档
- 优化代码结构
- 添加代码注释
- 准备项目发布

测试编写完成，让我们继续完善项目！
