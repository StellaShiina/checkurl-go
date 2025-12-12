package report

import (
	"fmt"
	"sort"

	"github.com/StellaShiina/checkurl-go/checker"
)

type Report struct {
	OK    []checker.CheckResult
	Fail  []checker.CheckResult
	ERROR []checker.CheckResult
}

func (r *Report) GenerateSummary(results map[string]checker.CheckResult) {
	for _, res := range results {
		if res.Error != nil {
			r.ERROR = append(r.ERROR, res)
		} else if res.StatusCode >= 200 && res.StatusCode < 300 {
			r.OK = append(r.OK, res)
		} else {
			r.Fail = append(r.Fail, res)
		}
	}
	r.sort2xxResults()
}

func (r *Report) sort2xxResults() {
	sort.Slice(r.OK, func(i, j int) bool {
		return r.OK[i].Latency < r.OK[j].Latency
	})
}

func (r *Report) PrintSummary() {
	fmt.Printf("OK: %d, Fail: %d, Error: %d\n", len(r.OK), len(r.Fail), len(r.ERROR))
}

func (r *Report) PrintOK() {
	if len(r.OK) == 0 {
		fmt.Println("No success...")
		return
	}
	fmt.Println("-----OK-----")
	for i, res := range r.OK {
		fmt.Printf("[%d] %s %s\n", i, res.URL, res.Latency.String())
	}
}

func (r *Report) PrintFail() {
	if len(r.Fail) == 0 {
		fmt.Println("Nothing failed...")
		return
	}
	fmt.Println("-----Fail-----")
	for i, res := range r.Fail {
		fmt.Printf("[%d] %s %s\n", i, res.URL, res.StatusText)
	}
}

func (r *Report) PrintError() {
	if len(r.ERROR) == 0 {
		fmt.Println("No error...")
		return
	}
	fmt.Println("-----Error-----")
	for i, res := range r.ERROR {
		fmt.Printf("[%d] %s %s\n", i, res.URL, res.Error)
	}
}
