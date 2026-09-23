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

	start := time.Now()
	ScheduleJobs(mockBenchParams, &mockWg, &mockChan)	

	mockWg.Wait()
	close(mockChan)

	elapsed := time.Since(start)
	elapsed_mili := time.Duration.Milliseconds(elapsed)

	var values []time.Duration = []time.Duration{}
	for v := range mockChan {
		values = append(values, v)	
	}
	
	if len(values) != mockBenchParams.reqCount {
		t.Errorf("Executed the wrong amount of jobs: Executed %v Needed %v", len(values), mockBenchParams.reqCount)	
	}

	expectedMinimumTime := mockBenchParams.reqCount * mockBenchParams.delay

	if int(elapsed_mili) < expectedMinimumTime {
		t.Errorf("Jobs took longer than needed: Expected Minimum: %v Got: %v", expectedMinimumTime, elapsed_mili)	
	}
}

