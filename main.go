package main

import (
	"flag"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"
)

/*
	==========
	BENCHKITTY
	==========
*/

type statusCount map[int]string

type benchmark struct {
	avg int
	p50 int
	p95 int
	p99 int
	rps int
	tps int	
	statusCount statusCount
}

type benchParams struct {
	url      string
	method   string
	delay    int
	reqCount int
	out      string
}

var methods = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
}

func isValidHttpMethod(s string) bool {
	return slices.Contains(methods, s)
}

func ScheduleJobs(p benchParams, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < p.reqCount {
		time.Sleep(time.Duration(p.delay) * time.Millisecond)

		go func() {
			start := time.Now()

			_, err := http.Get(p.url)
			if err != nil {
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

func calculateBenchmark(params benchParams, values []int) benchmark {
	len := len(values)
	sorted := append([]int(nil), values...)
	slices.Sort(sorted)

	total_elapsed := 0
	for _, elapsed := range sorted {
		total_elapsed += int(elapsed)
	}

	avg := total_elapsed / params.reqCount
	middle := len / 2
	p50 := sorted[middle]
	p95 := sorted[int(float32(len)*0.95)]
	p99 := sorted[int(float32(len)*0.99)]

	return benchmark{
		avg: avg,
		p50: p50,
		p95: p95,
		p99: p99,
	}
}

func main() {
	reqCount := flag.Int("r", 100, "request amount")
	delay := flag.Int("d", 1000, "delay between every request call in ms")
	out := flag.String("o", "", "output path for the benchmark xlsx (defaults to terminal)")
	method := flag.String("t", "GET", "http method used on the url")

	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		fmt.Println("Missing url")
		flag.Usage()
		return
	}

	if !isValidHttpMethod(*method) {
		fmt.Printf("Invalid value %v for command param -t\n", *method)
		return
	}

	params := benchParams{
		url:      url,
		method:   *method,
		delay:    *delay,
		reqCount: *reqCount,
		out:      *out,
	}

	var wg sync.WaitGroup
	wg.Add(params.reqCount)

	timeChan := make(chan time.Duration, params.reqCount)

	ScheduleJobs(params, &wg, &timeChan)

	wg.Wait()
	close(timeChan)

	latencies := []int{}
	for elapsed := range timeChan {
		latencies = append(latencies, int(time.Duration.Milliseconds(elapsed)))
	}

	benchmark := calculateBenchmark(params, latencies)

	if params.out != "" {
		writeXLSX(params, benchmark, latencies)
		return
	}

	fmt.Printf("\n%+v\n", benchmark)
}
