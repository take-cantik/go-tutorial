package main

import "fmt"

// Step 1: Go の基本
// 実行: make step1  または  go run ./steps/01-hello
func main() {
	name := "Go Tutorial"
	fmt.Printf("Hello, %s!\n", name)

	tasks := []string{"環境構築", "HTTP", "MySQL", "レイヤー分離", "クリーンアーキテクチャ"}
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task)
	}

	status := describeStatus("doing")
	fmt.Println("現在のステータス:", status)
}

func describeStatus(s string) string {
	switch s {
	case "todo":
		return "未着手"
	case "doing":
		return "進行中"
	case "done":
		return "完了"
	default:
		return "不明"
	}
}
