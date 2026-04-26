package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultOldName = "go-base"

func main() {
	oldName := flag.String("from", defaultOldName, "原模块名称")
	newName := flag.String("name", "", "新模块名称")
	flag.Parse()

	if *newName == "" {
		fmt.Println("用法: go run cmd/rename/main.go -from=原名称 -name=新名称")
		fmt.Println("示例: go run cmd/rename/main.go -from=go-base-server -name=server")
		os.Exit(1)
	}

	if *oldName == "" {
		fmt.Println("原模块名称不能为空")
		os.Exit(1)
	}

	if *oldName == *newName {
		fmt.Println("原名称和新名称相同，无需处理")
		os.Exit(0)
	}

	rootDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("准备将模块名从 [%s] 改为 [%s]\n", *oldName, *newName)

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

		if !shouldProcessFile(info.Name()) {
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
		newContent := strings.ReplaceAll(oldContent, *oldName, *newName)

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

func shouldProcessFile(name string) bool {
	if name == "go.mod" || name == "go.sum" || name == "package.json" || name == "package-lock.json" ||
		name == "Dockerfile" || name == "docker-compose.yml" || name == ".gitignore" {
		return true
	}

	textExtensions := []string{
		".go",
		".tpl",
		".md",
		".yaml",
		".yml",
		".json",
		".sql",
		".xml",
		".txt",
		".ts",
		".vue",
		".js",
		".mjs",
		".cjs",
		".sh",
		".ps1",
		".env",
		".toml",
	}

	for _, ext := range textExtensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}

	return false
}
