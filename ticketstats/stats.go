package ticketstats

import (
	"fmt"
	"log"
	"time"

	"github.com/montanaflynn/stats"
)

// Stats groups a average work statistics result for a ticket list.
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

// Timerages groups the time rages used for ticket statistics
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
func OldBugs(issues []*Issue, config Config) []*Issue {
	oldBugs := make([]*Issue, 0)

	bugs := OlderThanOneMonth(FilterByType(issues, config.Types.Bug))
	for _, issue := range bugs {
		if !issue.IsResolved() {
			oldBugs = append(oldBugs, issue)
		}
	}

	return oldBugs
}

// ResolutionTime calculates resolution time statistics for a ticket list.
func ResolutionTime(issues []*Issue) Stats {
	times := make([]float64, 0)

	for _, issue := range issues {
		// ignore tickets with no time bookings
		if issue.TimeSpend > 0.0 {
			times = append(times, float64(issue.TimeSpend))
		}
	}

	// Case: no tickets with time bookings.
	if len(times) == 0 {
		return Stats{
			Mean:   Work(0),
			Median: Work(0),
			Count:  0,
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

// ResultionTimesByType calculates the resolution time statistics for each
// ticket type in the given list.
func ResultionTimesByType(issues []*Issue) map[string]TimeRanges {
	result := make(map[string]TimeRanges)

	for _, t := range Types(issues) {
		typeIssues := FilterByType(issues, t)

		var tr TimeRanges

		typeIssuesRange := ClosedLastYear(typeIssues)
		tr.Year = ResolutionTime(typeIssuesRange)

		if tr.Year.Count == 0 {
			// no booked work hours
			continue
		}

		typeIssuesRange = ClosedLastQuarter(typeIssuesRange)
		tr.Quarter = ResolutionTime(typeIssuesRange)

		typeIssuesRange = ClosedLastMonth(typeIssuesRange)
		tr.Month = ResolutionTime(typeIssuesRange)

		typeIssuesRange = ClosedLastWeek(typeIssuesRange)
		tr.Week = ResolutionTime(typeIssuesRange)

		result[t] = tr
	}

	return result
}

// WorkAfter sums all work done after a given start date.
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
