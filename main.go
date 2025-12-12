package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/StellaShiina/checkurl-go/checker"
	"github.com/StellaShiina/checkurl-go/fileio"
	"github.com/StellaShiina/checkurl-go/report"
)

func main() {
	fmt.Println("Hello! This is URL checker! A project for practising concurrency!")

	inFile := flag.String("i", "", "input file")
	inDir := flag.String("d", "", "input dir")
	outDir := flag.String("o", "", "output dir")
	numWorkers := flag.Int("n", 20, "max concurrency")

	flag.Parse()

	urls := make([]string, 0, 200)

	if *inFile != "" {
		err := fileio.ReadUrlsFromDir(*inFile, &urls)
		if err != nil {
			panic(err)
		}
	}

	if *inDir != "" {
		files, err := os.ReadDir(*inDir)

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == "example.txt" {
				continue
			}
			f := filepath.Join(*inDir, file.Name())
			err := fileio.ReadUrlsFromDir(f, &urls)
			if err != nil {
				panic(err)
			}
		}
	}

	numJobs := len(urls)

	jobs := make(chan string, numJobs)
	results := make(chan checker.CheckResult, numJobs)
	var wg sync.WaitGroup

	for w := range *numWorkers {
		wg.Add(1)
		go checker.Worker(w, jobs, results, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}

	urls = nil

	close(jobs)

	wg.Wait()
	close(results)

	finalResults := make(map[string]checker.CheckResult)
	for result := range results {
		finalResults[result.URL] = result
	}

	var report = &report.Report{}

	report.GenerateSummary(finalResults)
	report.PrintSummary()
	report.PrintOK()
	// report.PrintFail()
	report.PrintError()

	if *outDir == "" {
		*outDir = filepath.Join("data", "results", "results.txt")
	}
	err := fileio.WriteResultsToFile(*outDir, report.OK, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mission accomplished!")
}
