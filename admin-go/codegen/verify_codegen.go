//go:build ignore

package main

import (
	"fmt"
	"os"

	"gbaseadmin/codegen/internal/verifytemplates"
)

// codegen 离线模板验证 — 覆盖当前协作约定中的 codegen 验收场景
// 运行: cd admin-go/codegen && go run verify_codegen.go

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("[verify_codegen] 获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	summary, err := verifytemplates.RunCLI(rootDir)
	if err != nil {
		fmt.Printf("[verify_codegen] 执行失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n========== 结果 ==========\n")
	fmt.Printf("总检查: %d, 失败: %d\n", summary.TotalChecks, summary.TotalErrors)
	if summary.TotalErrors > 0 {
		os.Exit(1)
	}
	fmt.Println("全部通过!")
}
