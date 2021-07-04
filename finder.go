package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	location *string
)

func init() {
	location = flag.String("location", "", "location to search")
}

type QueryList struct {
	QueryList []string `json:"query_list"`
}

func loadConfig() []string {
	file, _ := ioutil.ReadFile("config.json")

	data := QueryList{}

	_ = json.Unmarshal([]byte(file), &data)
	return data.QueryList
}

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
func searchEngine(files []string, query_list []string, result chan string) {

	for _, file := range files {
		for _, name := range query_list {
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
	flag.Parse()
	root := *location
	if len(root) == 0 {
		fmt.Println("Usage: finder.exe -location path/to/location")
		// flag.PrintDefaults()
		os.Exit(1)
	}

	query_list := loadConfig()
	result := make(chan string)

	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	go searchEngine(files, query_list, result)
	fileWriter(result)

}
