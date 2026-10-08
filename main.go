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
		panic("Params were unable to be parsed")
	}

	var wg sync.WaitGroup
	wg.Add(params.ReqCount)

	jobChan := make(chan benchkitty.JobInfo, params.ReqCount)

	start := time.Now()
	benchkitty.ScheduleJobs(params, &wg, &jobChan)

	wg.Wait()
	close(jobChan)
	duration := time.Since(start)

	jobs := []benchkitty.JobInfo{}
	for job := range jobChan {
		jobs = append(jobs, job)
	}

	result := benchkitty.CalculateBenchmark(params, jobs, duration)

	if params.Out != "" {
		benchkitty.XLSX(params, result, jobs)
		return
	}

	fmt.Printf("\n%+v\n", result)
}
