package main

import (
	"encoding/json"
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
	median int
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
	job_count := 17
	delay := 500

	var wg sync.WaitGroup
	wg.Add(job_count)

	time_chan := make(chan time.Duration, job_count)
	count := 0

	for count < job_count {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		go func() {
			start := time.Now()

			fmt.Println("Job Started")

			_, err := http.Get(url)
			if err != nil {
				panic("Request falhou")
			}

			elapsed := time.Since(start)
			
			wg.Done()
			
			time_chan <- elapsed
		}()

		count++
	}

	wg.Wait()
	close(time_chan)

	total_elapsed := 0
	for elapsed := range time_chan {
		fmt.Println(elapsed)
		elapsed_mili := time.Duration.Milliseconds(elapsed)
		total_elapsed += int(elapsed_mili)
	}	

	median := total_elapsed / job_count
	
	stats := stats{
		median: median,
	}
	
	json_stats, err := json.Marshal(stats)
	if err != nil{
		panic("Unable to parse stats struct")
	}
	
	fmt.Println(json_stats)
}
