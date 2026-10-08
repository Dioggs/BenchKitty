package benchkitty

import (
	"fmt"
	"slices"
	"strings"
)

type reportRow struct {
	label string
	value string
}

type reportSection struct {
	title string
	rows  []reportRow
}

func humanizeBytesPerSecond(bps int) string {
	units := []string{"B/s", "KB/s", "MB/s", "GB/s", "TB/s"}
	value := float64(bps)
	unit := 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d %s", bps, units[unit])
	}
	return fmt.Sprintf("%.2f %s", value, units[unit])
}

func buildReportSections(result Result) []reportSection {
	sections := []reportSection{
		{
			title: "Latency",
			rows: []reportRow{
				{"Avg", fmt.Sprintf("%d ms", result.Avg)},
				{"P50", fmt.Sprintf("%d ms", result.P50)},
				{"P95", fmt.Sprintf("%d ms", result.P95)},
				{"P99", fmt.Sprintf("%d ms", result.P99)},
			},
		},
		{
			title: "Throughput",
			rows: []reportRow{
				{"Requests/sec", fmt.Sprintf("%d req/s", result.RPS)},
				{"Bandwidth", humanizeBytesPerSecond(result.TPS)},
			},
		},
	}

	if len(result.StatusCount) > 0 {
		statuses := make([]int, 0, len(result.StatusCount))
		total := 0
		for status, count := range result.StatusCount {
			statuses = append(statuses, status)
			total += count
		}
		slices.Sort(statuses)

		rows := make([]reportRow, 0, len(statuses))
		for _, status := range statuses {
			count := result.StatusCount[status]
			rows = append(rows, reportRow{
				label: fmt.Sprintf("%d", status),
				value: fmt.Sprintf("%d (%.0f%%)", count, float64(count)/float64(total)*100),
			})
		}
		sections = append(sections, reportSection{title: "Status codes", rows: rows})
	}

	return sections
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func PrintResult(result Result) {
	title := "BenchKitty Benchmark Report"
	sections := buildReportSections(result)

	labelWidth, valueWidth := 0, 0
	for _, section := range sections {
		for _, row := range section.rows {
			if len(row.label) > labelWidth {
				labelWidth = len(row.label)
			}
			if len(row.value) > valueWidth {
				valueWidth = len(row.value)
			}
		}
	}

	inner := 2 + labelWidth + 2 + valueWidth
	if len(title)+2 > inner {
		inner = len(title) + 2
	}

	border := strings.Repeat("─", inner)
	titlePad := (inner - len(title)) / 2
	titleLine := strings.Repeat(" ", titlePad) + title
	titleLine = padRight(titleLine, inner)

	fmt.Println()
	fmt.Println("╭" + border + "╮")
	fmt.Println("│" + titleLine + "│")
	for i, section := range sections {
		fmt.Println("├" + border + "┤")
		if i > 0 {
			fmt.Println("│" + padRight("", inner) + "│")
		}
		fmt.Println("│ " + padRight(section.title, inner-1) + "│")
		for _, row := range section.rows {
			line := "  " + padRight(row.label, labelWidth+2) + padLeft(row.value, valueWidth)
			fmt.Println("│ " + padRight(line, inner-1) + "│")
		}
	}
	fmt.Println("╰" + border + "╯")
	fmt.Println()
}
