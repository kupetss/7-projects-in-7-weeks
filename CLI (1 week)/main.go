package main

import (
	"bufio"
	"flag"
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
	var (
		linesFlag = flag.Bool("l", false, "")
		wordsFlag = flag.Bool("w", false, "")
		bytesFlag = flag.Bool("c", false, "")
		fileFlag  = flag.String("file", "", "")
	)

	flag.Parse()

	if !*linesFlag && !*wordsFlag && !*bytesFlag {
		*linesFlag, *wordsFlag, *bytesFlag = true, true, true
	}

	path := *fileFlag
	if path == "" && flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	if path == "" {
		fmt.Println("Не указан путь к файлу или папке")
		return
	}

	info, _ := os.Stat(path)

	var results []Stats

	if info.IsDir() {
		results = folder(path)
	} else {
		stats := count(path)
		results = append(results, stats)
	}

	ShowResults(results, *linesFlag, *wordsFlag, *bytesFlag)
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

func ShowResults(results []Stats, showLines, showWords, showBytes bool) {
	maxLen := 0
	for _, item := range results {
		if len(item.Name) > maxLen {
			maxLen = len(item.Name)
		}
	}

	header := fmt.Sprintf("%-*s", maxLen, "Файл")
	if showLines {
		header += " | Строки"
	}
	if showWords {
		header += " | Слова"
	}
	if showBytes {
		header += " | Байты"
	}
	fmt.Println(header)

	fmt.Println(strings.Repeat("-", len(header)))

	for _, item := range results {
		line := fmt.Sprintf("%-*s", maxLen, item.Name)
		if showLines {
			line += fmt.Sprintf(" | %6d", item.Lines)
		}
		if showWords {
			line += fmt.Sprintf(" | %4d", item.Words)
		}
		if showBytes {
			line += fmt.Sprintf(" | %6d", item.Size)
		}
		fmt.Println(line)
	}
}
