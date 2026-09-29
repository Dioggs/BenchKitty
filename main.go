package main

import (
	"fmt"
	"sync"
	"time"

	benchkitty "BenchKitty/internal"
)

func main() {
	params, ok := benchkitty.Parse()
	if !ok {
		return
	}

	var wg sync.WaitGroup
	wg.Add(params.ReqCount)

	timeChan := make(chan time.Duration, params.ReqCount)

	benchkitty.ScheduleJobs(params, &wg, &timeChan)

	wg.Wait()
	close(timeChan)

	latencies := []int{}
	for elapsed := range timeChan {
		latencies = append(latencies, int(time.Duration.Milliseconds(elapsed)))
	}

	result := benchkitty.CalculateBenchmark(params, latencies)

	if params.Out != "" {
		benchkitty.XLSX(params, result, latencies)
		return
	}

	fmt.Printf("\n%+v\n", result)
}
