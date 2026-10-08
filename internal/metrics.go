package benchkitty

import (
	"slices"
	"time"
)

type StatusCount map[int]int

type Result struct {
	Avg         int
	P50         int
	P95         int
	P99         int
	RPS         int
	TPS         int
	StatusCount StatusCount
}

func CalculateBenchmark(params Config, jobs []JobInfo, duration time.Duration) Result {
	values := jobLatencies(jobs)

	count := len(values)
	sorted := append([]int(nil), values...)
	slices.Sort(sorted)

	total_elapsed := 0
	for _, elapsed := range sorted {
		total_elapsed += int(elapsed)
	}

	avg := total_elapsed / params.ReqCount
	middle := count / 2
	p50 := sorted[middle]
	p95 := sorted[int(float32(count)*0.95)]
	p99 := sorted[int(float32(count)*0.99)]

	statusCount := StatusCount{}
	totalThroughput := int64(0)
	for _, job := range jobs {
		statusCount[job.Status]++
		totalThroughput += job.Throughput
	}

	rps := 0
	tps := 0
	if seconds := duration.Seconds(); seconds > 0 {
		rps = int(float64(len(jobs)) / seconds)
		tps = int(float64(totalThroughput) / seconds)
	}

	return Result{
		Avg:         avg,
		P50:         p50,
		P95:         p95,
		P99:         p99,
		RPS:         rps,
		TPS:         tps,
		StatusCount: statusCount,
	}
}
