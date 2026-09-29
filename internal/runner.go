package benchkitty

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func ScheduleJobs(p Config, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < p.ReqCount {
		time.Sleep(time.Duration(p.Delay) * time.Millisecond)

		go func() {
			start := time.Now()

			if _, err := http.Get(p.URL); err != nil {
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
