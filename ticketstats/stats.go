package ticketstats

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
)

type Stats struct {
	Mean   Work
	Median Work
	Count  int
}

func (stats Stats) ToString() string {
	return fmt.Sprintf("mean: %s, median: %s, count: %d",
		formatWork(stats.Mean),
		formatWork(stats.Median),
		stats.Count)
}

type TimeRanges struct {
	Week    Stats
	Month   Stats
	Quarter Stats
	Year    Stats
}

func (tr TimeRanges) ToString() string {
	str := ""

	stats := tr.Week
	str += fmt.Sprintf("last week:    mean: %.15s, median: %.15s, %d issues\n",
		formatWork(stats.Mean),
		formatWork(stats.Median),
		stats.Count)

	stats = tr.Month
	str += fmt.Sprintf("last month:   mean: %.15s, median: %.15s, %d issues\n",
		formatWork(stats.Mean),
		formatWork(stats.Median),
		stats.Count)

	stats = tr.Quarter
	str += fmt.Sprintf("last quarter: mean: %.15s, median: %.15s, %d issues\n",
		formatWork(stats.Mean),
		formatWork(stats.Median),
		stats.Count)

	stats = tr.Year
	str += fmt.Sprintf("last year:    mean: %.15s, median: %.15s, %d issues\n",
		formatWork(stats.Mean),
		formatWork(stats.Median),
		stats.Count)

	return str
}

// OldBugs finds all unresolved bugs older than one month.
func OldBugs(issues []*Issue) []*Issue {
	oldBugs := make([]*Issue, 0)

	for _, issue := range OlderThanOneMonth(FilterByType(issues, "Bug")) {
		if !issue.IsResolved() {
			oldBugs = append(oldBugs, issue)
		}
	}

	return oldBugs
}

func ResolutionTime(issues []*Issue, startDate time.Time) Stats {
	times := make([]float64, 0)

	for _, issue := range issues {
		if issue.TimeSpend > 0.0 {
			times = append(times, float64(issue.TimeSpend))
		}
	}

	mean, err := stats.Mean(times)
	if err != nil {
		log.Println("ERROR: mean of resolution time", err)
	}

	median, err := stats.Median(times)
	if err != nil {
		log.Println("ERROR: median of resolution time", err)
	}

	return Stats{
		Mean:   Work(mean),
		Median: Work(median),
		Count:  len(times),
	}
}

func ResultionTimesByType(issues []*Issue) map[string]TimeRanges {
	result := make(map[string]TimeRanges)

	for _, t := range containedTypes(issues) {
		typeIssues := FilterByType(issues, t)

		var tr TimeRanges

		typeIssuesRange := ClosedLastYear(typeIssues)
		tr.Year = ResolutionTime(typeIssuesRange, time.Now().AddDate(-1, 0, 0))

		if tr.Year.Count == 0 {
			// no booked work hours
			continue
		}

		typeIssuesRange = ClosedLastQuarter(typeIssuesRange)
		tr.Quarter = ResolutionTime(typeIssuesRange, time.Now().AddDate(0, -3, 0))

		typeIssuesRange = ClosedLastMonth(typeIssuesRange)
		tr.Month = ResolutionTime(typeIssuesRange, time.Now().AddDate(0, -1, 0))

		typeIssuesRange = ClosedLastWeek(typeIssuesRange)
		tr.Week = ResolutionTime(typeIssuesRange, time.Now().AddDate(0, 0, -7))

		result[t] = tr
	}

	return result
}

func containedTypes(issues []*Issue) []string {
	typesString := ""
	types := make([]string, 0)

	for _, issue := range issues {
		t := issue.Type
		if !strings.Contains(typesString, t) {
			typesString += t
			types = append(types, t)
		}
	}

	return types
}

func WorkAfter(issues []*Issue, start time.Time) Work {
	var work Work
	for _, issue := range issues {
		for _, log := range issue.LogWorks {
			if log.Date.After(start) {
				work += log.Hours
			}
		}
	}
	return work
}
