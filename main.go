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

type benchmark struct {
	avg int
}

type benchParams struct {
	url        string
	method     string
	req_count  int
	delay      int
	output_dir string
}

const reqCount int = 17
var cmdParams []string = []string{
	"-r",
	"-d",
	"-o",
	"-t",
}

func isValidCmdParam(cmdParam string) bool {
	return true
}

func isValidCmdValue(cmdValue string) bool {
	return true	
}

func buildParams(args []string) benchParams {
	benchParams := benchParams{}
	
	for _, arg := range args {
		/*
			se for cmdParam, pega o proximo cmdValue e valida ele		
			se for valido, pula pro proximo param
			se não for valido, panic
			
			se for baseValue, valida ele e segue (permitimos só url por enquanto)
		*/
		fmt.Println(arg)
	}

	return benchParams
}

func scheduleJobs(url string, delay int, wg *sync.WaitGroup, ch *chan time.Duration) {
	count := 0

	for count < reqCount {
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
	wg.Add(reqCount)

	time_chan := make(chan time.Duration, reqCount)

	scheduleJobs(url, delay, &wg, &time_chan)

	wg.Wait()
	close(time_chan)

	total_elapsed := 0
	for elapsed := range time_chan {
		elapsed_mili := time.Duration.Milliseconds(elapsed)
		total_elapsed += int(elapsed_mili)
	}

	avg := total_elapsed / reqCount 

	benchmark := benchmark{
		avg: avg,
	}

	switch output {
	case "term":
		fmt.Printf("\n%+v\n", benchmark)
	case "csv":
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
	default:
		panic("Invalid output format")
	}

}
