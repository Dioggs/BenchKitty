package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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

			fmt.Println("Job Done")

			*ch <- elapsed
			wg.Done()
		}()

		count++
	}
}

// TODO: fix measurements
func calculateBenchmark(params benchParams, timeChan chan time.Duration) benchmark {
	var values []int

    for value := range timeChan{
        values = append(values , int(time.Duration.Milliseconds(value)))
	}
	
	len := len(values)
	slices.Sort(values)
	
	total_elapsed := 0
	for _, elapsed := range values {
		total_elapsed += int(elapsed)
	}

	avg := total_elapsed / params.reqCount
	middle := len / 2
	p50 := values[middle]

	benchmark := benchmark{
		avg: avg,
		p50: p50,
	}

	return benchmark
}

func writeCSV(params benchParams, benchmark benchmark) {
	file, err := os.Create(filepath.Join(params.out, "benchmark.csv"))
	if err != nil {
		panic("failed to create file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	var data [][]string = [][]string{{"AVG"}, {strconv.Itoa(benchmark.avg)}}

	if err := writer.WriteAll(data); err != nil {
		panic("failed to write data")
	}
}

func main() {
	reqCount := flag.Int("r", 100, "request amount")
	delay := flag.Int("d", 1000, "delay between every request call in ms")
	out := flag.String("o", "", "output path for the benchmark csv (defaults to terminal)")
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

	benchmark := calculateBenchmark(params, timeChan)

	if params.out != "" {
		writeCSV(params, benchmark)	
		return
	}

	fmt.Printf("\n%+v\n", benchmark)
}
