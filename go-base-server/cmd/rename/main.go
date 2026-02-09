package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const oldModuleName = "ecoBreed-server"

func main() {
	newName := flag.String("name", "", "新模块名称")
	flag.Parse()

	if *newName == "" {
		fmt.Println("用法: go run cmd/rename/main.go -name=新模块名")
		fmt.Println("示例: go run cmd/rename/main.go -name=my-project")
		os.Exit(1)
	}

	rootDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("准备将模块名从 [%s] 改为 [%s]\n", oldModuleName, *newName)

	count := 0

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和非目标文件
		if info.IsDir() {
			// 跳过 .git 和 vendor 目录
			if info.Name() == ".git" || info.Name() == "vendor" || info.Name() == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 go.mod、.go 和 .tpl 文件
		if info.Name() != "go.mod" && !strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), ".tpl") {
			return nil
		}

		// 跳过当前工具自身
		if strings.Contains(path, "cmd"+string(os.PathSeparator)+"rename") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		oldContent := string(content)
		newContent := strings.ReplaceAll(oldContent, oldModuleName, *newName)

		if oldContent != newContent {
			if err := os.WriteFile(path, []byte(newContent), info.Mode()); err != nil {
				return err
			}
			relPath, _ := filepath.Rel(rootDir, path)
			fmt.Printf("✓ %s\n", relPath)
			count++
		}

		return nil
	})

	if err != nil {
		fmt.Printf("处理失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n完成! 共更新 %d 个文件\n", count)
	fmt.Println("请执行: go mod tidy")
}
