package benchkitty

import (
	"testing"
	"time"
)

func TestCalculateBenchmarkAggregatesStatusAndThroughput(t *testing.T) {
	jobs := []JobInfo{
		{Status: 200, Throughput: 100, Elapsed: 10 * time.Millisecond},
		{Status: 200, Throughput: 300, Elapsed: 20 * time.Millisecond},
		{Status: 500, Throughput: 0, Elapsed: 30 * time.Millisecond},
	}

	result := CalculateBenchmark(Config{ReqCount: 3}, jobs, time.Second)

	if result.StatusCount[200] != 2 {
		t.Errorf("expected 2 responses with status 200, got %v", result.StatusCount[200])
	}
	if result.StatusCount[500] != 1 {
		t.Errorf("expected 1 response with status 500, got %v", result.StatusCount[500])
	}
	if result.TPS != 400 {
		t.Errorf("expected throughput 400 bytes/s, got %v", result.TPS)
	}
	if result.RPS != 3 {
		t.Errorf("expected 3 rps, got %v", result.RPS)
	}
}
