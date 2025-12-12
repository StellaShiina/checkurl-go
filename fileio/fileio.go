package fileio

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/StellaShiina/checkurl-go/checker"
)

func ReadUrlsFromDir(dir string, urls *[]string) error {
	fmt.Printf("Start to read from %s...\n", dir)
	fileHandler, err := os.Open(dir)
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
			fmt.Printf("File %s, line %d is not valid\n", dir, count)
			continue
		}
		*urls = append(*urls, line)
	}
	fmt.Printf("Finish reading from %s...\n", dir)
	return nil
}

func WriteResultsToFile(dir string, results []checker.CheckResult, onlyURL bool) error {
	fmt.Printf("Start to write results to %s, onlyURL: %t\n", dir, onlyURL)
	f, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, res := range results {
		line := res.URL + "\n"
		if !onlyURL {
			line = res.URL + res.StatusText + res.Latency.String() + "\n"
		}
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Finish writing results to %s\n", dir)
	return nil
}
