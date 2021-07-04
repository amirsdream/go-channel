package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func stringFinder(files []string, result chan string) {
	names := [3]string{"es2015", "node", "test"}
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

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), name) {
			result <- file + " " + name
		}
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
		if matched, err := filepath.Match("*.json", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			*files = append(*files, path)
		}
		return nil
	}
}

func main() {
	var files []string
	result := make(chan string)

	root := "D:\\Code"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	go stringFinder(files, result)

	for v := range result {
		if v != "none" {
			fmt.Println(v)
		}
	}

}
