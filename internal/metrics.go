package benchkitty

import "slices"

type StatusCount map[int]string

type Result struct {
	Avg         int
	P50         int
	P95         int
	P99         int
	RPS         int
	TPS         int
	StatusCount StatusCount
}

func CalculateBenchmark(params Config, values []int) Result {
	len := len(values)
	sorted := append([]int(nil), values...)
	slices.Sort(sorted)

	total_elapsed := 0
	for _, elapsed := range sorted {
		total_elapsed += int(elapsed)
	}

	avg := total_elapsed / params.ReqCount
	middle := len / 2
	p50 := sorted[middle]
	p95 := sorted[int(float32(len)*0.95)]
	p99 := sorted[int(float32(len)*0.99)]

	return Result{
		Avg: avg,
		P50: p50,
		P95: p95,
		P99: p99,
	}
}
