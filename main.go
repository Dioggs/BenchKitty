package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
	==========
	BENCHKITTY
	==========
*/

type stats struct {
	avg int
}

const REQ_COUNT int = 17

func schedule_jobs(url string, delay int, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < REQ_COUNT {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		go func() {
			start := time.Now()

			_, err := http.Get(url)
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
	delay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Unable to parse delay")
	}

	output := os.Args[2]
	url := os.Args[3]

	var wg sync.WaitGroup
	wg.Add(REQ_COUNT)

	time_chan := make(chan time.Duration, REQ_COUNT)

	schedule_jobs(url, delay, &wg, &time_chan)

	wg.Wait()
	close(time_chan)

	total_elapsed := 0
	for elapsed := range time_chan {
		elapsed_mili := time.Duration.Milliseconds(elapsed)
		total_elapsed += int(elapsed_mili)
	}

	avg := total_elapsed / REQ_COUNT

	stats := stats{
		avg: avg,
	}

	switch output {
	case "term":
		fmt.Printf("\n%+v\n", stats)
	case "csv":
		file, err := os.Create("benchmark.csv")
		if err != nil {
			panic("failed to create file")
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		var data [][]string = [][]string{{"AVG"}, {strconv.Itoa(stats.avg)}}

		if err := writer.WriteAll(data); err != nil {
			panic("failed to write data")
		}
	default:
		panic("Invalid output format")
	}

}
