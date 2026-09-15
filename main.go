package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

/*
	BENCHKITTY
*/

type stats struct {
	avg int
}

const JOB_COUNT int = 17
const DELAY int = 500

func schedule_jobs(url string, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < JOB_COUNT {
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
	/*
		recebo url como parametro

		posso dizer quantidade requests
		posso determinar body
		posso dizer quantidade de jobs paralelos
		posso salvar métrica em CSV
		posso dizer o tipo do request
		posso determinar timeout

		retorno métrica
		existe a possibilidade de uma TUI bonitinha
	*/

	url := os.Args[1]

	var wg sync.WaitGroup
	wg.Add(JOB_COUNT)

	time_chan := make(chan time.Duration, JOB_COUNT)
	
	schedule_jobs(url, &wg, &time_chan)

	wg.Wait()
	close(time_chan)

	total_elapsed := 0
	for elapsed := range time_chan {
		elapsed_mili := time.Duration.Milliseconds(elapsed)
		total_elapsed += int(elapsed_mili)
	}

	avg := total_elapsed / JOB_COUNT

	stats := stats{
		avg: avg,
	}

	fmt.Printf("\n%+v\n", stats)
}
