package main

import (
	"fmt"
	"net/http"
	"os"
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
const DELAY int = 500

func schedule_jobs(url string, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < REQ_COUNT {
		time.Sleep(time.Duration(DELAY) * time.Millisecond)

		go func() {
			start := time.Now()

			_, err := http.Get(url)
			if err != nil {
				panic("Request Failed")
			}

			elapsed := time.Since(start)

			wg.Done()
			fmt.Println("Job Done")

			*ch <- elapsed
		}()

		count++
	}
}

func main() {
	url := os.Args[1]

	var wg sync.WaitGroup
	wg.Add(REQ_COUNT)

	time_chan := make(chan time.Duration, REQ_COUNT)

	schedule_jobs(url, &wg, &time_chan)

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

	fmt.Printf("\n%+v\n", stats)
}
