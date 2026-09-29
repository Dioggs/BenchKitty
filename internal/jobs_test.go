package benchkitty

import (
	"sync"
	"testing"
	"time"
)

func TestSchedulingOfJobs(t *testing.T) {
	mockBenchParams := Config{
		URL:      "https://jsonplaceholder.typicode.com/todos/1",
		ReqCount: 10,
		Delay:    100,
		Method:   "GET",
		Out:      "",
	}
	mockChan := make(chan time.Duration, mockBenchParams.ReqCount)
	var mockWg sync.WaitGroup
	mockWg.Add(mockBenchParams.ReqCount)

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

	if len(values) != mockBenchParams.ReqCount {
		t.Errorf("Executed the wrong amount of jobs: Executed %v Needed %v", len(values), mockBenchParams.ReqCount)
	}

	expectedMinimumTime := mockBenchParams.ReqCount * mockBenchParams.Delay

	if int(elapsed_mili) < expectedMinimumTime {
		t.Errorf("Jobs took longer than needed: Expected Minimum: %v Got: %v", expectedMinimumTime, elapsed_mili)
	}
}
