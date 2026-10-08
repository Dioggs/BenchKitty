package benchkitty

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

func ScheduleJobs(p Config, wg *sync.WaitGroup, ch *chan time.Duration) {
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
				panic("Request Failed")
			}

			status := resp.StatusCode
			resp.Body.Close()

			elapsed := time.Since(start)

			if p.Pretty {
				fmtOutput(p.URL, p.Method, status, elapsed)
			}

			*ch <- elapsed
			wg.Done()
		}()

		count++
	}
}
