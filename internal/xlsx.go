package benchkitty

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/xuri/excelize/v2"
)

func XLSX(params Config, result Result, jobs []JobInfo) {
	file := excelize.NewFile()
	defer file.Close()

	if err := file.SetSheetName("Sheet1", "Summary"); err != nil {
		panic("failed to create summary sheet")
	}
	if _, err := file.NewSheet("Latency"); err != nil {
		panic("failed to create latency sheet")
	}
	if _, err := file.NewSheet("Status"); err != nil {
		panic("failed to create status sheet")
	}

	file.SetCellValue("Summary", "A1", "Metric")
	file.SetCellValue("Summary", "B1", "Value")
	file.SetCellValue("Summary", "A2", "Avg")
	file.SetCellValue("Summary", "B2", result.Avg)
	file.SetCellValue("Summary", "A3", "P50")
	file.SetCellValue("Summary", "B3", result.P50)
	file.SetCellValue("Summary", "A4", "P95")
	file.SetCellValue("Summary", "B4", result.P95)
	file.SetCellValue("Summary", "A5", "P99")
	file.SetCellValue("Summary", "B5", result.P99)
	file.SetCellValue("Summary", "A6", "RPS")
	file.SetCellValue("Summary", "B6", result.RPS)
	file.SetCellValue("Summary", "A7", "Throughput (bytes/s)")
	file.SetCellValue("Summary", "B7", result.TPS)

	file.SetCellValue("Latency", "A1", "Request")
	file.SetCellValue("Latency", "B1", "Latency (ms)")
	file.SetCellValue("Latency", "C1", "Sorted (ms)")

	latencies := jobLatencies(jobs)
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

	statuses := make([]int, 0, len(result.StatusCount))
	for status := range result.StatusCount {
		statuses = append(statuses, status)
	}
	slices.Sort(statuses)

	file.SetCellValue("Status", "A1", "Status")
	file.SetCellValue("Status", "B1", "Count")
	for i, status := range statuses {
		row := i + 2
		file.SetCellValue("Status", fmt.Sprintf("A%d", row), status)
		file.SetCellValue("Status", fmt.Sprintf("B%d", row), result.StatusCount[status])
	}

	lastStatus := len(statuses) + 1
	if err := file.AddChart("Status", "D2", &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{{
			Name:       "Status!$B$1",
			Categories: fmt.Sprintf("Status!$A$2:$A$%d", lastStatus),
			Values:     fmt.Sprintf("Status!$B$2:$B$%d", lastStatus),
		}},
		Title: excelize.ChartTitle{
			Paragraph: []excelize.RichTextRun{{Text: "Status code distribution"}},
		},
	}); err != nil {
		panic("failed to add status chart")
	}

	if err := file.SaveAs(filepath.Join(params.Out, "benchmark.xlsx")); err != nil {
		panic("failed to save xlsx")
	}
}
