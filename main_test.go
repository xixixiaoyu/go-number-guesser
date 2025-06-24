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
