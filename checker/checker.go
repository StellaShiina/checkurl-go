package checker

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var defaultClient = &http.Client{Timeout: 5 * time.Second}

const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

type CheckResult struct {
	URL        string
	StatusCode int
	StatusText string
	Latency    time.Duration
	Error      error
}

func check(c *http.Client, url string) (*http.Response, time.Duration, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")
	start := time.Now()
	resp, err := c.Do(req)
	latency := time.Since(start)
	return resp, latency, err
}

func Worker(id int, jobs <-chan string, result chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		resp, latency, err := check(defaultClient, url)
		if err != nil {
			result <- CheckResult{
				URL:        url,
				StatusCode: 0,
				StatusText: fmt.Sprintf("Network Error (%s)", err.Error()),
				Latency:    latency,
				Error:      err,
			}
		} else {
			result <- CheckResult{
				URL:        url,
				StatusCode: resp.StatusCode,
				StatusText: resp.Status,
				Latency:    latency,
				Error:      nil,
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}
