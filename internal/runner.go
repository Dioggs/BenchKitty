package benchkitty

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type JobInfo struct {
	Status     int
	Throughput int64
	Elapsed    time.Duration
}

func jobLatencies(jobs []JobInfo) []int {
	values := make([]int, 0, len(jobs))
	for _, job := range jobs {
		values = append(values, int(time.Duration.Milliseconds(job.Elapsed)))
	}
	return values
}

func buildRequest(p Config) (*http.Request, error) {
	if p.Body == "" {
		return http.NewRequest(p.Method, p.URL, nil)
	}

	req, err := http.NewRequest(p.Method, p.URL, strings.NewReader(p.Body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func fmtOutput(url string, method string, status int, elapsed time.Duration) {
	fmt.Println()
	fmt.Printf("URL: %v\n", url)
	fmt.Printf("Method: %v\n", method)
	fmt.Printf("Status: %v\n", status)
	fmt.Printf("Elapsed: %vms\n", time.Duration.Milliseconds(elapsed))
	fmt.Println()
}

func ScheduleJobs(p Config, wg *sync.WaitGroup, ch *chan JobInfo) {
	count := 0

	for count < p.ReqCount {
		time.Sleep(time.Duration(p.Delay) * time.Millisecond)

		go func() {
			start := time.Now()

			req, err := buildRequest(p)
			if err != nil {
				panic("Unable to build request")
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic("Unable to execute request")
			}

			status := resp.StatusCode
			throughput, _ := io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			elapsed := time.Since(start)

			if p.Pretty {
				fmtOutput(p.URL, p.Method, status, elapsed)
			}

			*ch <- JobInfo{
				Elapsed:    elapsed,
				Status:     status,
				Throughput: throughput,
			}
			wg.Done()
		}()

		count++
	}
}
