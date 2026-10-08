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

			if _, err := http.DefaultClient.Do(req); err != nil {
				panic("Request Failed")
			}

			elapsed := time.Since(start)

			fmt.Printf("Job Done = %vms\n", time.Duration.Milliseconds(elapsed))

			*ch <- elapsed
			wg.Done()
		}()

		count++
	}
}
