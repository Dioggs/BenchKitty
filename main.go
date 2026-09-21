package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
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

type benchmark struct {
	avg int
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

func isValidPath(s string) bool {
	return fs.ValidPath(s)
}

func isValidHttpMethod(s string) bool {
	return slices.Contains(methods, s)
}

func scheduleJobs(p benchParams, wg *sync.WaitGroup, ch *chan time.Duration) {
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

	if *out != "" && !isValidPath(*out) {
		fmt.Printf("Invalid value %v for command param -o\n", *out)
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

	time_chan := make(chan time.Duration, params.reqCount)

	scheduleJobs(params, &wg, &time_chan)

	wg.Wait()
	close(time_chan)

	total_elapsed := 0
	for elapsed := range time_chan {
		elapsed_mili := time.Duration.Milliseconds(elapsed)
		total_elapsed += int(elapsed_mili)
	}

	avg := total_elapsed / params.reqCount

	benchmark := benchmark{
		avg: avg,
	}

	if params.out != "" {
		// TODO: handle "params.out" as a relative path
		file, err := os.Create("benchmark.csv")
		if err != nil {
			panic("failed to create file")
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		var data [][]string = [][]string{{"AVG"}, {strconv.Itoa(benchmark.avg)}}

		if err := writer.WriteAll(data); err != nil {
			panic("failed to write data")
		}
		
		return
	}
	
	fmt.Printf("\n%+v\n", benchmark)
}
