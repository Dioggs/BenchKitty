package benchkitty

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestMethodAndBodyAreSent(t *testing.T) {
	type request struct {
		method string
		body   string
	}
	captured := make(chan request, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured <- request{method: r.Method, body: string(body)}
	}))
	defer server.Close()

	mockBenchParams := Config{
		URL:      server.URL,
		Method:   "POST",
		Body:     `{"hello":"world"}`,
		ReqCount: 1,
		Delay:    10,
	}

	mockChan := make(chan JobInfo, mockBenchParams.ReqCount)
	var mockWg sync.WaitGroup
	mockWg.Add(mockBenchParams.ReqCount)

	ScheduleJobs(mockBenchParams, &mockWg, &mockChan)

	mockWg.Wait()
	close(mockChan)

	got := <-captured
	if got.method != mockBenchParams.Method {
		t.Errorf("Sent the wrong method: Got %v Needed %v", got.method, mockBenchParams.Method)
	}
	if got.body != mockBenchParams.Body {
		t.Errorf("Sent the wrong body: Got %v Needed %v", got.body, mockBenchParams.Body)
	}
}

func TestSchedulingOfJobs(t *testing.T) {
	mockBenchParams := Config{
		URL:      "https://jsonplaceholder.typicode.com/todos/1",
		ReqCount: 10,
		Delay:    100,
		Method:   "GET",
		Out:      "",
	}
	mockChan := make(chan JobInfo, mockBenchParams.ReqCount)
	var mockWg sync.WaitGroup
	mockWg.Add(mockBenchParams.ReqCount)

	start := time.Now()
	ScheduleJobs(mockBenchParams, &mockWg, &mockChan)

	mockWg.Wait()
	close(mockChan)

	elapsed := time.Since(start)
	elapsed_mili := time.Duration.Milliseconds(elapsed)

	var jobs []JobInfo = []JobInfo{}
	for job := range mockChan {
		jobs = append(jobs, job)
	}

	if len(jobs) != mockBenchParams.ReqCount {
		t.Errorf("Executed the wrong amount of jobs: Executed %v Needed %v", len(jobs), mockBenchParams.ReqCount)
	}

	expectedMinimumTime := mockBenchParams.ReqCount * mockBenchParams.Delay

	if int(elapsed_mili) < expectedMinimumTime {
		t.Errorf("Jobs took longer than needed: Expected Minimum: %v Got: %v", expectedMinimumTime, elapsed_mili)
	}
}
