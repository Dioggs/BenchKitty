package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
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

type benchParams map[string]string

const reqCount int = 17

func isValidInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isValidPath(s string) bool {
	return fs.ValidPath(s)	
}

func isValidHttpMethod(s string) bool {
	methods := []string {
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
	}
	
	return slices.Contains(methods, s)
}
 
func getErrMsg(cmdParam string, cmdValue string) string {
	return fmt.Sprintf("Invalid value %v for command param %v", cmdValue, cmdParam)
}

func validateCmdParam(cmdParam string, cmdValue string) error {
	switch cmdParam {
	case "-d":
		res := isValidInt(cmdValue)
		if !res {
			return errors.New(getErrMsg(cmdParam, cmdValue))
		}
	case "-r": 
		res := isValidInt(cmdValue)
		if !res {
			return errors.New(getErrMsg(cmdParam, cmdValue))
		}
	case "-o":
		res := isValidPath(cmdValue)
		if !res {
			return errors.New(getErrMsg(cmdParam, cmdValue))
		}
	case "-t":
		res := isValidHttpMethod(cmdValue)
		if !res {
			return errors.New(getErrMsg(cmdParam, cmdValue))
		}
	default: 
		return errors.New("Unsupported command parameter")
	}	
	
	return nil
}

func buildBenchParams(args []string) benchParams {
	benchParams := make(benchParams)
	len := len(args)

	for i:=0; i<len; i+=2{
		arg := args[i]
		if strings.Contains(arg, "-") {
			if i + 1 >= len {
				panic("Missing value for cmd param " + arg)
			}	
			
			val := args[i + 1]
			
			err := validateCmdParam(arg, val)
			if err != nil {
				panic(err)
			}
			benchParams[arg] = val
		}	
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

func main(){
	args := os.Args[1:]
	params := buildBenchParams(args)
	fmt.Println(params)
}

func main2() {
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
