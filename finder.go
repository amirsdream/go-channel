package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func fileWriter(result chan string) {
	path := "test.csv"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for v := range result {
		if v != "none" {
			_, _ = datawriter.WriteString(v + "\n")
		}
	}

	datawriter.Flush()
	file.Close()
}
func searchEngine(files []string, names []string, result chan string) {

	for _, file := range files {
		for _, name := range names {
			stringProcessor(file, name, result)
		}
	}
	close(result)
}

func stringProcessor(file string, name string, result chan string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), name) {
			result <- file + ", " + name + ", " + strconv.Itoa(line)
		}
		line++
	}

	result <- "none"
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
			// log.Fatal(err)
		}
		re := regexp.MustCompile(`\.json$|\.txt$`)
		if re.Match([]byte(filepath.Base(path))) {
			*files = append(*files, path)
		}
		return nil
	}
}

func main() {
	var files []string
	names := []string{"es2015", "test"}
	result := make(chan string)

	root := "D:\\Code"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	go searchEngine(files, names, result)
	fileWriter(result)

}
