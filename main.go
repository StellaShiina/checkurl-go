package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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

	const numWorkers = 20
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

	var numOK, numFail, numError = 0, 0, 0

	res := make([]checker.CheckResult, 0, 100)

	for k, v := range finalResults {
		if v.StatusCode < 200 && v.StatusCode > 0 || v.StatusCode >= 300 {
			numFail++
			fmt.Println(k, v.Latency, v.StatusCode, v.StatusText)
		} else if v.StatusCode == 0 {
			numError++
			fmt.Println(k, v.Error)
		} else {
			numOK++
			res = append(res, v)
			fmt.Println(k, v.Latency, v.StatusCode)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Latency < res[j].Latency
	})

	fmt.Printf("OK: %d, Fail: %d, Error: %d\n", numOK, numFail, numError)

	fmt.Println("OK URLs:")

	for i, r := range res {
		fmt.Println(i, r.URL, r.Latency)
	}
}
