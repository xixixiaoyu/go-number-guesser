# 第六章：游戏逻辑

## 🎮 游戏逻辑深度分析

在前面的章节中，我们已经实现了基本的游戏功能。现在让我们深入分析游戏逻辑的实现细节，理解每个组件是如何协同工作的。

## 🧠 核心算法分析

### 随机数生成算法

```go
// 在 NewGame() 函数中
rand.Seed(time.Now().UnixNano())
targetNumber: rand.Intn(100) + 1
```

#### 算法原理

1. **伪随机数生成器**：
   - Go 使用线性同余生成器（LCG）
   - 公式：`X(n+1) = (a * X(n) + c) mod m`
   - 需要种子值来初始化序列

2. **种子设置**：
   ```go
   rand.Seed(time.Now().UnixNano())
   ```
   - 使用当前时间的纳秒值作为种子
   - 确保每次运行程序都有不同的随机序列
   - `UnixNano()` 返回自 1970 年以来的纳秒数

3. **范围映射**：
   ```go
   rand.Intn(100) + 1
   ```
   - `rand.Intn(100)` 生成 [0, 100) 范围的整数
   - `+1` 将范围调整为 [1, 100]

#### 随机性质量分析

```go
// 测试随机数分布的代码示例
func testRandomDistribution() {
    counts := make(map[int]int)
    
    for i := 0; i < 10000; i++ {
        num := rand.Intn(100) + 1
        counts[num]++
    }
    
    // 理论上每个数字应该出现约 100 次
    for num, count := range counts {
        fmt.Printf("数字 %d 出现了 %d 次\n", num, count)
    }
}
```

### 比较算法

```go
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

#### 算法特点

1. **时间复杂度**：O(1) - 常数时间
2. **空间复杂度**：O(1) - 常数空间
3. **返回值设计**：
   - `0`：相等（猜对了）
   - `1`：guess > target（猜大了）
   - `-1`：guess < target（猜小了）

#### 为什么这样设计？

1. **简洁性**：三种情况用三个不同的数值表示
2. **可扩展性**：未来可以根据返回值添加更多逻辑
3. **测试友好**：返回值明确，易于编写测试用例

## 🔄 游戏状态管理

### 状态变量分析

```go
type Game struct {
    targetNumber int           // 目标数字状态
    attempts     int           // 统计状态
    scanner      *bufio.Scanner // I/O 状态
}
```

#### 状态生命周期

1. **初始化阶段**：
   ```go
   func NewGame() *Game {
       return &Game{
           targetNumber: rand.Intn(100) + 1, // 设置目标
           attempts:     0,                   // 重置计数
           scanner:      bufio.NewScanner(os.Stdin), // 初始化输入
       }
   }
   ```

2. **游戏进行阶段**：
   ```go
   func (g *Game) Start() {
       for {
           // 获取输入
           // 更新状态：g.attempts++
           // 检查结果
           // 如果猜对，退出循环
       }
   }
   ```

3. **游戏结束阶段**：
   - 显示最终统计信息
   - 实例被垃圾回收器回收

### 状态一致性保证

1. **原子操作**：
   ```go
   g.attempts++ // 只有在有效输入后才增加
   ```

2. **状态隔离**：
   - 每个游戏实例独立
   - 不同游戏之间不会相互影响

3. **不变量维护**：
   - `targetNumber` 在游戏过程中不变
   - `attempts` 只能递增
   - `scanner` 保持有效状态

## 🎯 游戏流程控制

### 主循环设计

```go
func main() {
    // 显示标题
    for {
        game := NewGame()    // 创建新游戏
        game.Start()         // 开始游戏
        
        if !askContinue() {  // 询问是否继续
            break            // 退出主循环
        }
        
        // 显示新游戏标题
    }
    // 显示结束信息
}
```

#### 流程控制分析

1. **外层循环**：控制多轮游戏
2. **内层循环**：控制单轮游戏（在 `Start()` 方法中）
3. **条件退出**：用户选择不继续时退出

### 单轮游戏流程

```go
func (g *Game) Start() {
    // 显示欢迎信息
    
    for {
        // 1. 获取用户输入
        guess, err := g.getPlayerGuess()
        
        // 2. 错误处理
        if err != nil {
            fmt.Printf("输入错误：%v，请重新输入。\n", err)
            continue // 重新开始循环，不增加计数
        }
        
        // 3. 更新状态
        g.attempts++
        
        // 4. 检查结果
        result := g.checkGuess(guess)
        
        // 5. 根据结果采取行动
        if result == 0 {
            // 猜对了，显示成功信息并退出
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

#### 流程特点

1. **错误恢复**：输入错误不会终止游戏
2. **状态更新**：只有有效输入才更新计数
3. **及时反馈**：每次猜测都有即时反馈
4. **明确退出**：猜对后立即退出循环

## 🔍 输入验证逻辑

### 多层验证策略

```go
func (g *Game) getPlayerGuess() (int, error) {
    fmt.Print("请输入你的猜测：")
    
    // 第一层：读取验证
    if !g.scanner.Scan() {
        return 0, fmt.Errorf("读取输入失败")
    }
    
    // 第二层：空值验证
    input := strings.TrimSpace(g.scanner.Text())
    if input == "" {
        return 0, fmt.Errorf("输入不能为空")
    }
    
    // 第三层：格式验证
    guess, err := strconv.Atoi(input)
    if err != nil {
        return 0, fmt.Errorf("请输入一个有效的数字")
    }
    
    // 第四层：范围验证
    if guess < 1 || guess > 100 {
        return 0, fmt.Errorf("数字必须在 1-100 之间")
    }
    
    return guess, nil
}
```

#### 验证层次分析

1. **系统级验证**：确保能够读取输入
2. **格式级验证**：确保输入不为空
3. **类型级验证**：确保输入是数字
4. **业务级验证**：确保数字在有效范围内

#### 错误处理策略

1. **早期返回**：一旦发现错误立即返回
2. **具体错误信息**：每种错误都有明确的提示
3. **用户友好**：错误信息使用中文，易于理解

## 🎲 随机性与公平性

### 随机性测试

```go
// 测试随机数质量的函数
func TestRandomness(t *testing.T) {
    numbers := make(map[int]bool)
    
    // 生成多个游戏实例，检查目标数字的多样性
    for i := 0; i < 20; i++ {
        game := NewGame()
        numbers[game.targetNumber] = true
    }
    
    // 至少应该有几个不同的数字
    if len(numbers) < 3 {
        t.Errorf("随机数生成缺乏多样性，20次生成只有 %d 个不同数字", len(numbers))
    }
}
```

### 公平性保证

1. **每个数字等概率**：理论上每个数字被选中的概率相等
2. **独立性**：每轮游戏的目标数字相互独立
3. **不可预测性**：玩家无法预测下一个目标数字

## 📊 性能分析

### 时间复杂度分析

1. **NewGame()**：O(1) - 常数时间
2. **checkGuess()**：O(1) - 常数时间
3. **getPlayerGuess()**：O(1) - 常数时间（不考虑用户输入时间）
4. **整个游戏**：O(n) - n 为猜测次数

### 空间复杂度分析

1. **Game 结构体**：O(1) - 固定大小
2. **输入缓冲区**：O(1) - 固定大小
3. **整个程序**：O(1) - 常数空间

### 性能优化考虑

1. **避免不必要的内存分配**：
   ```go
   // 好的做法：复用 scanner
   scanner: bufio.NewScanner(os.Stdin)
   
   // 不好的做法：每次都创建新的 scanner
   ```

2. **减少字符串操作**：
   ```go
   // 只在必要时进行字符串处理
   input := strings.TrimSpace(g.scanner.Text())
   ```

## 🎯 算法优化思考

### 当前算法的优势

1. **简单直观**：逻辑清晰，易于理解
2. **高效执行**：所有操作都是 O(1) 时间复杂度
3. **内存友好**：空间使用量很小
4. **易于测试**：每个函数职责单一

### 可能的改进方向

1. **智能提示**：
   ```go
   // 可以添加更智能的提示
   if attempts > 10 {
       fmt.Println("提示：试试二分查找策略！")
   }
   ```

2. **难度调节**：
   ```go
   type Game struct {
       targetNumber int
       attempts     int
       scanner      *bufio.Scanner
       minRange     int // 新增：最小范围
       maxRange     int // 新增：最大范围
   }
   ```

3. **统计功能**：
   ```go
   type GameStats struct {
       totalGames    int
       totalAttempts int
       bestScore     int
   }
   ```

## 🧪 逻辑测试策略

### 单元测试设计

1. **边界条件测试**：
   - 输入 1 和 100
   - 目标数字为 1 和 100

2. **正常流程测试**：
   - 多次猜测直到成功
   - 验证计数器正确性

3. **异常情况测试**：
   - 各种无效输入
   - 系统错误模拟

## 🎯 本章总结

在这一章中，我们深入分析了游戏逻辑：

1. ✅ **随机数算法**：理解了伪随机数生成原理
2. ✅ **比较算法**：分析了简单高效的比较逻辑
3. ✅ **状态管理**：掌握了游戏状态的生命周期
4. ✅ **流程控制**：理解了多层循环的控制逻辑
5. ✅ **输入验证**：学习了多层验证策略
6. ✅ **性能分析**：了解了算法的时间和空间复杂度

### 关键收获

- **算法设计**：简单往往是最好的
- **错误处理**：多层验证确保程序健壮性
- **状态管理**：清晰的状态设计便于维护
- **性能考虑**：在简单性和性能之间找到平衡

## 🚀 下一步

在下一章中，我们将：
- 编写完整的测试套件
- 学习单元测试和基准测试
- 掌握测试驱动开发方法
- 确保代码质量和可靠性

游戏逻辑分析完成，让我们继续学习测试技术！
