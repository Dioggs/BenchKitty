package main

import (
	"sync"
	"testing"
	"time"
)

func TestSchedulingOfJobs(t *testing.T) {
	mockBenchParams := benchParams{
		url: "https://jsonplaceholder.typicode.com/todos/1",
		reqCount: 10,
		delay: 100,
		method: "GET",
		out: "",
	}
	mockChan := make(chan time.Duration, mockBenchParams.reqCount)
	var mockWg sync.WaitGroup
	mockWg.Add(mockBenchParams.reqCount)

	ScheduleJobs(mockBenchParams, &mockWg, &mockChan)	
	
	mockWg.Wait()
	close(mockChan)
	
	// did we execute the current amout of jobs?
	// did all of them return data?
	// did we take the correct amount of time?
	// 
}

