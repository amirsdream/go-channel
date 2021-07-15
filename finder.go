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

type JsonData struct {
	QueryList []string `json:"query_list"`
	FileTypes []string `json:"file_types"`
}

func loadConfig() JsonData {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("config file has an issue: %s", err)
	}

	data := JsonData{}

	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func fileWriter(result chan string) {
	path := "test.csv"

	for msg := range result {
		fmt.Println("received", msg)
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		l, err := f.WriteString(msg + "\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "bytes written successfully")
		err = f.Close()
	}
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
	// result <- "none"
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func visit(files *[]string, file_types []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
			// log.Fatal(err)
		}
		pattern := ""
		for _, t := range file_types {
			temp := fmt.Sprintf(`\.%s$|`, t)
			pattern += temp
		}
		pattern = pattern[:len(pattern)-1]
		re := regexp.MustCompile(pattern)
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
	data := loadConfig()
	result := make(chan string)
	if len(root) == 0 {
		fmt.Println("Usage: finder.exe -location path/to/location")
		os.Exit(1)
	}
	err := filepath.Walk(root, visit(&files, data.FileTypes))
	if err != nil {
		panic(err)
	}
	go searchEngine(files, data.QueryList, result)

	fileWriter(result)

}
