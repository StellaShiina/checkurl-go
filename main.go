package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/StellaShiina/checkurl-go/checker"
)

func read(file string, urls *[]string) error {
	fileHandler, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileHandler.Close()

	scanner := bufio.NewScanner(fileHandler)

	count := -1

	for scanner.Scan() {
		count++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		u, err := url.Parse(line)
		if err != nil || u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
			fmt.Printf("File %s, line %d is not valid\n", file, count)
			continue
		}
		*urls = append(*urls, line)
	}

	return nil
}

func main() {
	fmt.Println("Hello! This is URL checker! A project for practising concurrency!")

	const dataDir = "data"

	urls := make([]string, 0, 200)

	files, err := os.ReadDir(dataDir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f := filepath.Join(dataDir, file.Name())
		err := read(f, &urls)
		if err != nil {
			panic(err)
		}
	}

	numJobs := len(urls)

	jobs := make(chan string, numJobs)
	results := make(chan checker.CheckResult, numJobs)
	var wg sync.WaitGroup

	const numWorkers = 10
	for w := range numWorkers {
		wg.Add(1)
		go checker.Worker(w, jobs, results, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}

	close(jobs)

	wg.Wait()
	close(results)

	finalResults := make(map[string]checker.CheckResult)
	for result := range results {
		finalResults[result.URL] = result
	}

	for k, v := range finalResults {
		fmt.Println(k, v.StatusCode)
	}
}
