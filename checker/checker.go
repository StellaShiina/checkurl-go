package checker

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var defaultClient = &http.Client{Timeout: 5 * time.Second}

// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
// Accept-Language: en-US,en
// Accept-Encoding: gzip, deflate, br, zstd
// Connection: keep-alive

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0"
const Accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
const AcceptLanguage = "en-US,en"
const AcceptEncoding = "gzip, deflate, br, zstd"
const Connection = "keep-alive"

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
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Accept-Encoding", AcceptEncoding)
	req.Header.Set("Connection", Connection)
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
