package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Game 结构体用于管理游戏状态
type Game struct {
	targetNumber int            // 目标数字
	attempts     int            // 猜测次数
	scanner      *bufio.Scanner // 输入扫描器
}

// NewGame 创建新游戏实例
func NewGame() *Game {
	// 设置随机数种子
	return &Game{
		targetNumber: rand.Intn(100) + 1, // 生成 1-100 之间的随机数
		attempts:     0,
		scanner:      bufio.NewScanner(os.Stdin),
	}
}

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
