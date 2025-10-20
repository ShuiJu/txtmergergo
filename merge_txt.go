package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("使用方法: go run merge_txt.go <输入文件夹> <输出文件>")
		fmt.Println("示例: go run merge_txt.go ./input merged.txt")
		fmt.Println("示例: go run merge_txt.go . result.txt")
		return
	}

	inputDir := os.Args[1]
	outputFile := os.Args[2]

	// 读取文件夹中的所有文件
	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("读取文件夹失败: %v\n", err)
		return
	}

	// 筛选出.txt文件并获取完整路径
	var txtFiles []string
	for _, file := range files {
		if !file.IsDir() && isTxtFile(file.Name()) {
			fullPath := filepath.Join(inputDir, file.Name())
			txtFiles = append(txtFiles, fullPath)
		}
	}

	if len(txtFiles) == 0 {
		fmt.Println("文件夹中没有找到txt文件")
		return
	}

	// 按照文件名升序排序
	sort.Strings(txtFiles)

	fmt.Printf("找到 %d 个txt文件，按文件名升序排列:\n", len(txtFiles))
	for i, file := range txtFiles {
		fmt.Printf("%d. %s\n", i+1, filepath.Base(file))
	}

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("创建输出文件失败: %v\n", err)
		return
	}
	defer outFile.Close()

	// 按排序后的顺序处理文件
	for _, file := range txtFiles {
		if err := mergeFile(file, outFile); err != nil {
			fmt.Printf("处理文件 %s 失败: %v\n", file, err)
		} else {
			fmt.Printf("成功处理: %s\n", filepath.Base(file))
		}
	}

	fmt.Printf("\n合并完成！输出文件: %s\n", outputFile)
	fmt.Printf("总共处理了 %d 个文件\n", len(txtFiles))
}

func isTxtFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".txt" || ext == ".TXT"
}

func mergeFile(inputFile string, outFile *os.File) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// 写入文件名（只写基本文件名，不包含路径）
	fileName := filepath.Base(inputFile)
	_, err = fmt.Fprintf(outFile, "%s\n", fileName)
	if err != nil {
		return err
	}

	// 复制文件内容
	if _, err := io.Copy(outFile, inFile); err != nil {
		return err
	}

	// 写入分隔换行
	_, err = outFile.WriteString("\n")
	return err
}
