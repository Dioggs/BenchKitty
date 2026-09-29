package main

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/xuri/excelize/v2"
)

func writeXLSX(params benchParams, benchmark benchmark, latencies []int) {
	file := excelize.NewFile()
	defer file.Close()

	if err := file.SetSheetName("Sheet1", "Summary"); err != nil {
		panic("failed to create summary sheet")
	}
	if _, err := file.NewSheet("Latency"); err != nil {
		panic("failed to create latency sheet")
	}

	file.SetCellValue("Summary", "A1", "Metric")
	file.SetCellValue("Summary", "B1", "Value")
	file.SetCellValue("Summary", "A2", "Avg")
	file.SetCellValue("Summary", "B2", benchmark.avg)
	file.SetCellValue("Summary", "A3", "P50")
	file.SetCellValue("Summary", "B3", benchmark.p50)
	file.SetCellValue("Summary", "A4", "P95")
	file.SetCellValue("Summary", "B4", benchmark.p95)
	file.SetCellValue("Summary", "A5", "P99")
	file.SetCellValue("Summary", "B5", benchmark.p99)

	file.SetCellValue("Latency", "A1", "Request")
	file.SetCellValue("Latency", "B1", "Latency (ms)")
	file.SetCellValue("Latency", "C1", "Sorted (ms)")

	sorted := append([]int(nil), latencies...)
	slices.Sort(sorted)
	for i, latency := range latencies {
		row := i + 2
		file.SetCellValue("Latency", fmt.Sprintf("A%d", row), i+1)
		file.SetCellValue("Latency", fmt.Sprintf("B%d", row), latency)
		file.SetCellValue("Latency", fmt.Sprintf("C%d", row), sorted[i])
	}

	last := len(latencies) + 1

	if err := file.AddChart("Summary", "D2", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{{
			Name:       "Summary!$B$1",
			Categories: "Summary!$A$2:$A$5",
			Values:     "Summary!$B$2:$B$5",
		}},
		Title: excelize.ChartTitle{
			Paragraph: []excelize.RichTextRun{{Text: "Latency percentiles (ms)"}},
		},
	}); err != nil {
		panic("failed to add percentile chart")
	}

	if err := file.AddChart("Latency", "E2", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{{
			Name:       "Latency!$B$1",
			Categories: fmt.Sprintf("Latency!$A$2:$A$%d", last),
			Values:     fmt.Sprintf("Latency!$B$2:$B$%d", last),
		}},
		Title: excelize.ChartTitle{
			Paragraph: []excelize.RichTextRun{{Text: "Latency per request (ms)"}},
		},
	}); err != nil {
		panic("failed to add latency chart")
	}

	if err := file.AddChart("Latency", "E20", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{{
			Name:       "Latency!$C$1",
			Categories: fmt.Sprintf("Latency!$A$2:$A$%d", last),
			Values:     fmt.Sprintf("Latency!$C$2:$C$%d", last),
		}},
		Title: excelize.ChartTitle{
			Paragraph: []excelize.RichTextRun{{Text: "Sorted latency distribution (ms)"}},
		},
	}); err != nil {
		panic("failed to add distribution chart")
	}

	if err := file.SaveAs(filepath.Join(params.out, "benchmark.xlsx")); err != nil {
		panic("failed to save xlsx")
	}
}
