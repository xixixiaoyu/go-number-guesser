# 第十章：扩展功能

## 🚀 功能扩展思路

我们已经完成了一个功能完整的猜数字游戏，但这只是开始。在这一章中，我们将探讨如何扩展游戏功能，让它更加丰富和有趣。

## 🎯 难度级别系统

### 设计思路

添加不同的难度级别，让玩家可以选择挑战：

```go
// DifficultyLevel 表示游戏难度级别
type DifficultyLevel struct {
    Name     string // 难度名称
    MinRange int    // 最小数字
    MaxRange int    // 最大数字
    MaxTries int    // 最大尝试次数（0 表示无限制）
}

// 预定义的难度级别
var (
    EasyLevel   = DifficultyLevel{"简单", 1, 50, 0}
    NormalLevel = DifficultyLevel{"普通", 1, 100, 0}
    HardLevel   = DifficultyLevel{"困难", 1, 200, 15}
    ExpertLevel = DifficultyLevel{"专家", 1, 500, 20}
)
```

### 实现难度选择

```go
// 扩展 Game 结构体
type Game struct {
    targetNumber int
    attempts     int
    scanner      *bufio.Scanner
    difficulty   DifficultyLevel // 新增：难度级别
}

// selectDifficulty 让玩家选择难度级别
func selectDifficulty() DifficultyLevel {
    difficulties := []DifficultyLevel{
        EasyLevel, NormalLevel, HardLevel, ExpertLevel,
    }
    
    fmt.Println("请选择难度级别：")
    for i, diff := range difficulties {
        maxTries := "无限制"
        if diff.MaxTries > 0 {
            maxTries = fmt.Sprintf("%d次", diff.MaxTries)
        }
        fmt.Printf("%d. %s (%d-%d, 最多%s)\n", 
            i+1, diff.Name, diff.MinRange, diff.MaxRange, maxTries)
    }
    
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("请输入选择 (1-4): ")
        if !scanner.Scan() {
            continue
        }
        
        choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
        if err != nil || choice < 1 || choice > len(difficulties) {
            fmt.Println("无效选择，请重新输入。")
            continue
        }
        
        return difficulties[choice-1]
    }
}

// 修改 NewGame 函数
func NewGame() *Game {
    difficulty := selectDifficulty()
    rand.Seed(time.Now().UnixNano())
    
    return &Game{
        targetNumber: rand.Intn(difficulty.MaxRange-difficulty.MinRange+1) + difficulty.MinRange,
        attempts:     0,
        scanner:      bufio.NewScanner(os.Stdin),
        difficulty:   difficulty,
    }
}
```

## 📊 游戏统计系统

### 统计数据结构

```go
// GameStats 游戏统计信息
type GameStats struct {
    TotalGames     int     // 总游戏次数
    TotalAttempts  int     // 总猜测次数
    BestScore      int     // 最佳成绩（最少猜测次数）
    WorstScore     int     // 最差成绩（最多猜测次数）
    AverageScore   float64 // 平均猜测次数
    WinRate        float64 // 胜率（对于有限制次数的难度）
    DifficultyStats map[string]*DifficultyStats // 各难度统计
}

// DifficultyStats 特定难度的统计
type DifficultyStats struct {
    Games     int     // 该难度的游戏次数
    Wins      int     // 胜利次数
    BestScore int     // 该难度最佳成绩
    AvgScore  float64 // 该难度平均成绩
}

// 全局统计实例
var globalStats = &GameStats{
    DifficultyStats: make(map[string]*DifficultyStats),
}
```

### 统计功能实现

```go
// updateStats 更新游戏统计
func (stats *GameStats) updateStats(difficulty DifficultyLevel, attempts int, won bool) {
    stats.TotalGames++
    
    if won {
        stats.TotalAttempts += attempts
        
        if stats.BestScore == 0 || attempts < stats.BestScore {
            stats.BestScore = attempts
        }
        
        if attempts > stats.WorstScore {
            stats.WorstScore = attempts
        }
        
        stats.AverageScore = float64(stats.TotalAttempts) / float64(stats.TotalGames)
    }
    
    // 更新难度统计
    diffName := difficulty.Name
    if stats.DifficultyStats[diffName] == nil {
        stats.DifficultyStats[diffName] = &DifficultyStats{}
    }
    
    diffStats := stats.DifficultyStats[diffName]
    diffStats.Games++
    
    if won {
        diffStats.Wins++
        if diffStats.BestScore == 0 || attempts < diffStats.BestScore {
            diffStats.BestScore = attempts
        }
        diffStats.AvgScore = (diffStats.AvgScore*float64(diffStats.Wins-1) + float64(attempts)) / float64(diffStats.Wins)
    }
    
    // 计算胜率
    if difficulty.MaxTries > 0 {
        totalWins := 0
        totalGames := 0
        for _, ds := range stats.DifficultyStats {
            totalWins += ds.Wins
            totalGames += ds.Games
        }
        if totalGames > 0 {
            stats.WinRate = float64(totalWins) / float64(totalGames) * 100
        }
    }
}

// showStats 显示统计信息
func (stats *GameStats) showStats() {
    fmt.Println("\n" + strings.Repeat("=", 40))
    fmt.Println("📊 游戏统计")
    fmt.Println(strings.Repeat("=", 40))
    
    fmt.Printf("总游戏次数: %d\n", stats.TotalGames)
    if stats.TotalGames > 0 {
        fmt.Printf("最佳成绩: %d 次\n", stats.BestScore)
        fmt.Printf("最差成绩: %d 次\n", stats.WorstScore)
        fmt.Printf("平均成绩: %.1f 次\n", stats.AverageScore)
        
        if stats.WinRate > 0 {
            fmt.Printf("胜率: %.1f%%\n", stats.WinRate)
        }
        
        fmt.Println("\n各难度统计:")
        for diffName, diffStats := range stats.DifficultyStats {
            fmt.Printf("  %s: %d局, 胜利%d次", diffName, diffStats.Games, diffStats.Wins)
            if diffStats.BestScore > 0 {
                fmt.Printf(", 最佳%d次, 平均%.1f次", diffStats.BestScore, diffStats.AvgScore)
            }
            fmt.Println()
        }
    }
    fmt.Println(strings.Repeat("=", 40))
}
```

## 🎮 智能提示系统

### 提示策略

```go
// HintSystem 提示系统
type HintSystem struct {
    game     *Game
    hints    []string
    hintUsed int
}

// generateHints 生成智能提示
func (h *HintSystem) generateHints() {
    target := h.game.targetNumber
    minRange := h.game.difficulty.MinRange
    maxRange := h.game.difficulty.MaxRange
    
    h.hints = []string{
        fmt.Sprintf("数字在 %d 到 %d 之间", minRange, maxRange),
        fmt.Sprintf("数字%s", getParityHint(target)),
        fmt.Sprintf("数字%s", getDivisibilityHint(target)),
        fmt.Sprintf("数字在 %s", getRangeHint(target, minRange, maxRange)),
    }
}

// getParityHint 获取奇偶性提示
func getParityHint(number int) string {
    if number%2 == 0 {
        return "是偶数"
    }
    return "是奇数"
}

// getDivisibilityHint 获取整除性提示
func getDivisibilityHint(number int) string {
    divisors := []int{3, 5, 7, 11}
    for _, d := range divisors {
        if number%d == 0 {
            return fmt.Sprintf("能被 %d 整除", d)
        }
    }
    return "是质数或有其他因子"
}

// getRangeHint 获取范围提示
func getRangeHint(number, min, max int) string {
    quarter := (max - min) / 4
    if number <= min+quarter {
        return "前四分之一范围内"
    } else if number <= min+quarter*2 {
        return "第二个四分之一范围内"
    } else if number <= min+quarter*3 {
        return "第三个四分之一范围内"
    }
    return "最后四分之一范围内"
}

// getHint 获取提示
func (h *HintSystem) getHint() string {
    if h.hintUsed >= len(h.hints) {
        return "没有更多提示了！"
    }
    
    hint := h.hints[h.hintUsed]
    h.hintUsed++
    return hint
}
```

## 🏆 成就系统

### 成就定义

```go
// Achievement 成就
type Achievement struct {
    ID          string
    Name        string
    Description string
    Unlocked    bool
    UnlockTime  time.Time
}

// 预定义成就
var achievements = map[string]*Achievement{
    "first_win": {
        ID:          "first_win",
        Name:        "初出茅庐",
        Description: "完成第一局游戏",
    },
    "lucky_guess": {
        ID:          "lucky_guess",
        Name:        "运气爆棚",
        Description: "一次猜中答案",
    },
    "persistent": {
        ID:          "persistent",
        Name:        "坚持不懈",
        Description: "完成10局游戏",
    },
    "master": {
        ID:          "master",
        Name:        "猜数大师",
        Description: "平均猜测次数少于5次",
    },
    "expert_winner": {
        ID:          "expert_winner",
        Name:        "专家级玩家",
        Description: "在专家难度下获胜",
    },
}

// checkAchievements 检查成就
func checkAchievements(stats *GameStats, difficulty DifficultyLevel, attempts int) {
    // 首次获胜
    if stats.TotalGames == 1 && !achievements["first_win"].Unlocked {
        unlockAchievement("first_win")
    }
    
    // 一次猜中
    if attempts == 1 && !achievements["lucky_guess"].Unlocked {
        unlockAchievement("lucky_guess")
    }
    
    // 坚持不懈
    if stats.TotalGames >= 10 && !achievements["persistent"].Unlocked {
        unlockAchievement("persistent")
    }
    
    // 猜数大师
    if stats.AverageScore < 5 && stats.TotalGames >= 5 && !achievements["master"].Unlocked {
        unlockAchievement("master")
    }
    
    // 专家级玩家
    if difficulty.Name == "专家" && !achievements["expert_winner"].Unlocked {
        unlockAchievement("expert_winner")
    }
}

// unlockAchievement 解锁成就
func unlockAchievement(id string) {
    achievement := achievements[id]
    if achievement != nil && !achievement.Unlocked {
        achievement.Unlocked = true
        achievement.UnlockTime = time.Now()
        
        fmt.Println("\n🏆 成就解锁！")
        fmt.Printf("【%s】%s\n", achievement.Name, achievement.Description)
    }
}
```

## 💾 数据持久化

### 保存和加载数据

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
)

// GameData 游戏数据
type GameData struct {
    Stats        *GameStats               `json:"stats"`
    Achievements map[string]*Achievement  `json:"achievements"`
}

// getDataFilePath 获取数据文件路径
func getDataFilePath() string {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".guess-game-data.json")
}

// saveGameData 保存游戏数据
func saveGameData() error {
    data := &GameData{
        Stats:        globalStats,
        Achievements: achievements,
    }
    
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(getDataFilePath(), jsonData, 0644)
}

// loadGameData 加载游戏数据
func loadGameData() error {
    filePath := getDataFilePath()
    
    // 如果文件不存在，使用默认数据
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil
    }
    
    jsonData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    var data GameData
    if err := json.Unmarshal(jsonData, &data); err != nil {
        return err
    }
    
    if data.Stats != nil {
        globalStats = data.Stats
    }
    
    if data.Achievements != nil {
        for id, achievement := range data.Achievements {
            if achievements[id] != nil {
                achievements[id].Unlocked = achievement.Unlocked
                achievements[id].UnlockTime = achievement.UnlockTime
            }
        }
    }
    
    return nil
}
```

## 🎨 界面美化

### 彩色输出

```go
// 颜色常量
const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[37m"
)

// colorPrint 彩色打印
func colorPrint(color, text string) {
    fmt.Printf("%s%s%s", color, text, ColorReset)
}

// 使用示例
func displayColorfulTitle() {
    colorPrint(ColorCyan, strings.Repeat("=", 50)+"\n")
    colorPrint(ColorYellow, "🎯 Go 语言猜数字游戏\n")
    colorPrint(ColorCyan, strings.Repeat("=", 50)+"\n")
}

func displayResult(result int) {
    switch result {
    case 0:
        colorPrint(ColorGreen, "🎉 恭喜你！猜对了！\n")
    case 1:
        colorPrint(ColorRed, "📈 太大了！请再试一次。\n")
    case -1:
        colorPrint(ColorBlue, "📉 太小了！请再试一次。\n")
    }
}
```

## 🌐 多语言支持

### 国际化框架

```go
// Language 语言配置
type Language struct {
    Code     string
    Name     string
    Messages map[string]string
}

var languages = map[string]*Language{
    "zh": {
        Code: "zh",
        Name: "中文",
        Messages: map[string]string{
            "welcome":     "欢迎来到猜数字游戏！",
            "guess_prompt": "请输入你的猜测：",
            "too_high":    "太大了！请再试一次。",
            "too_low":     "太小了！请再试一次。",
            "correct":     "恭喜你！猜对了！",
            "attempts":    "你总共猜了 %d 次。",
        },
    },
    "en": {
        Code: "en",
        Name: "English",
        Messages: map[string]string{
            "welcome":      "Welcome to the Number Guessing Game!",
            "guess_prompt": "Enter your guess: ",
            "too_high":     "Too high! Try again.",
            "too_low":      "Too low! Try again.",
            "correct":      "Congratulations! You got it!",
            "attempts":     "You made %d attempts.",
        },
    },
}

var currentLanguage = languages["zh"] // 默认中文

// getMessage 获取本地化消息
func getMessage(key string, args ...interface{}) string {
    if msg, exists := currentLanguage.Messages[key]; exists {
        if len(args) > 0 {
            return fmt.Sprintf(msg, args...)
        }
        return msg
    }
    return key // 如果找不到，返回键名
}
```

## 🎯 扩展功能总结

我们探讨了多种扩展功能：

1. ✅ **难度级别系统**：不同范围和限制的挑战
2. ✅ **游戏统计系统**：详细的游戏数据分析
3. ✅ **智能提示系统**：帮助玩家的提示功能
4. ✅ **成就系统**：增加游戏趣味性的成就
5. ✅ **数据持久化**：保存游戏进度和统计
6. ✅ **界面美化**：彩色输出和更好的视觉效果
7. ✅ **多语言支持**：国际化功能

## 🚀 实现优先级建议

### 第一阶段（核心扩展）
1. 难度级别系统
2. 基本统计功能
3. 数据持久化

### 第二阶段（体验优化）
1. 智能提示系统
2. 界面美化
3. 成就系统

### 第三阶段（高级功能）
1. 多语言支持
2. 网络对战模式
3. 图形界面版本

## 📚 学习路径建议

### Go 语言进阶
1. **并发编程**：goroutine 和 channel
2. **网络编程**：HTTP 服务器和客户端
3. **数据库操作**：SQL 和 NoSQL 数据库
4. **微服务架构**：gRPC 和服务发现

### 项目实践
1. **Web 应用**：使用 Gin 或 Echo 框架
2. **API 服务**：RESTful API 设计
3. **命令行工具**：使用 Cobra 库
4. **系统工具**：文件处理、日志分析等

### 开源贡献
1. 参与 Go 社区项目
2. 发布自己的 Go 包
3. 编写技术博客
4. 分享学习经验

## 🎯 结语

通过这个猜数字游戏项目，我们学习了：

- **Go 语言基础**：语法、标准库、工具链
- **软件设计**：架构设计、模式应用
- **工程实践**：测试、文档、部署
- **项目管理**：版本控制、自动化构建

这只是 Go 语言学习的开始。继续探索，不断实践，你将成为一名优秀的 Go 开发者！

**Happy Coding! 🚀**
