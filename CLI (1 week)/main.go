package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Stats struct {
	Name  string
	Lines int
	Words int
	Size  int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Не указан путь к файлу или папке")
		fmt.Println("Пример: go run main.go ./test.txt")
		return
	}

	path := os.Args[1]

	info, _ := os.Stat(path)

	var results []Stats

	if info.IsDir() {
		results = folder(path)
	} else {
		stats := count(path)
		results = append(results, stats)
	}

	ShowResults(results)
}

func folder(fPath string) []Stats {
	var results []Stats
	extentions := []string{".txt", ".go", ".doc", ".docx", ".html", ".css"}
	filepath.Walk(fPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			for _, ext := range extentions {
				if strings.HasSuffix(currentPath, ext) {
					stats := count(currentPath)
					results = append(results, stats)
				}
			}
		}
		return nil
	})
	return results
}

func count(path string) Stats {
	file, _ := os.Open(path)
	defer file.Close()

	stats := Stats{Name: path}

	linesCount := 0
	lineScanner := bufio.NewScanner(file)
	for lineScanner.Scan() {
		linesCount++
	}
	stats.Lines = linesCount

	file.Seek(0, 0)

	wordsCount := 0
	wordsScanner := bufio.NewScanner(file)
	wordsScanner.Split(bufio.ScanWords)
	for wordsScanner.Scan() {
		wordsCount++
	}
	stats.Words = wordsCount

	info, _ := file.Stat()
	stats.Size = info.Size()

	return stats
}

func ShowResults(results []Stats) {
	maxLen := 0
	for _, item := range results {
		if len(item.Name) > maxLen {
			maxLen = len(item.Name)
		}
	}

	fmt.Printf("%-*s | Строки | Слова | Байты\n", maxLen, "Файл")
	fmt.Println(strings.Repeat("-", maxLen+30))
	for _, item := range results {
		fmt.Printf("%-*s | %6d | %4d | %6d\n",
			maxLen,
			item.Name,
			item.Lines,
			item.Words,
			item.Size)
	}
}
