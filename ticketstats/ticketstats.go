// ticketstats generates statistics and further report for a set of csv exported
// Jira tickets.
package ticketstats

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type TicketStats struct {
	config    Config
	jiraBase  string
	issues    []*Issue
	active    []*Issue
	report    Report
	ignoreOld bool
}

// Evaluate generates a full report for the exported tickets
func Evaluate(path string,
	project string,
	component string,
	jiraBase string,
	splitByComponent bool) {

	config := DefaultConfig()

	// read issues form csv
	issues := Parse(path, config)

	if project != "" {
		issues = FilterByProject(issues, project)
	}
	if component != "" {
		issues = FilterByComponent(issues, component)
		splitByComponent = false
	}

	ClusterIssues(issues)
	PrintClusters(issues)

	ts := TicketStats{
		config:   config,
		jiraBase: jiraBase,
		issues:   issues,
		report:   NewReport(),
	}
	ts.report.Component = component
	ts.report.Date = time.Now().Format("2006-01-02")
	ts.ignoreOld = true
	ts.generateReport()

	if splitByComponent {
		for _, component := range Components(issues) {
			issues = FilterByComponent(issues, component)

			ts = TicketStats{
				config:   DefaultConfig(),
				jiraBase: jiraBase,
				issues:   FilterByComponent(issues, component),
				report:   NewReport(),
			}
			ts.report.Component = component
			ts.report.Date = time.Now().Format("2006-01-02")
			ts.ignoreOld = true
			ts.generateReport()
		}
	}
}

func (ts *TicketStats) generateReport() {
	// Reduce to active tickets
	ts.active = ActiveTickets(ts.issues)
	log.Println("INFO:", len(ts.active), "active tickets.")

	ts.sanitize()
	ts.oldBugs()
	ts.bugs()
	ts.features()
	ts.improvements()
	ts.other()
	ts.resources()

	ts.report.Render()
}

func (ts *TicketStats) sanitize() {
	// Check tickets for issues
	result := Sanitize(ts.issues, ts.ignoreOld)
	ts.report.Warnings = result.ToWarnings(ts.jiraBase)
	if ts.report.Warnings.Count > 0 {
		ts.report.HasWarnings = true
	}
}

func (ts *TicketStats) oldBugs() {
	oldBugs := OldBugs(ts.active)

	filterStates := []string{"Verification", "Acceptace", "Integration"}

	oldBugs = Filter(oldBugs, func(issue *Issue) bool {
		keep := true
		for _, status := range filterStates {
			if issue.Status == status {
				keep = false
				break
			}
		}
		return keep
	})

	OrderByCreated(oldBugs)
	for _, bug := range oldBugs {
		ts.report.OldBugs = append(ts.report.OldBugs, bug.ToReportIssue(ts.jiraBase))
	}
	log.Println("INFO:", len(oldBugs), "old bug tickets.")
}

func (ts *TicketStats) bugs() {
	bugs := FilterByType(ts.issues, "Bug")
	openBugs := OpenTickets(bugs)

	filterStates := []string{"Verification", "Acceptace", "Integration"}

	openBugs = Filter(openBugs, func(issue *Issue) bool {
		keep := true
		for _, status := range filterStates {
			if issue.Status == status {
				keep = false
				break
			}
		}
		return keep
	})

	ts.report.Bugs.Count = len(openBugs)

	ts.report.Bugs.Week.Created = len(CreatedLastWeek(bugs))
	ts.report.Bugs.Week.Resolved = len(ClosedLastWeek(bugs))
	ts.report.Bugs.Week.Diff = ts.report.Bugs.Week.Created - ts.report.Bugs.Week.Resolved

	ts.report.Bugs.Month.Created = len(CreatedLastMonth(bugs))
	ts.report.Bugs.Month.Resolved = len(ClosedLastMonth(bugs))
	ts.report.Bugs.Month.Diff = ts.report.Bugs.Month.Created - ts.report.Bugs.Month.Resolved

	versions := FixVersions(openBugs)
	securityLevels := SecurityLevels(openBugs)
	sort.Slice(versions, func(i, j int) bool {
		return strings.Compare(versions[i], versions[j]) > 0
	})

	ts.report.Bugs.BugCounts.Versions = versions

	for _, security := range securityLevels {
		sbs := FilterBySecurityLevel(openBugs, security)
		values := make([]string, 0)
		values = append(values, security)
		sum := 0
		for _, version := range versions {
			bs := FilterByFixVersion(sbs, version)

			stat := NewReportBugStats()
			stat.Count = len(bs)
			sum += stat.Count

			if stat.Count > 0 {
				values = append(values, fmt.Sprintf("%d", stat.Count))

				stat.Version = version
				stat.Security = security

				OrderByStatus(bs)
				OrderByPriority(bs)
				for _, b := range bs {
					stat.Bugs = append(stat.Bugs, b.ToReportIssue(ts.jiraBase))
				}

				ts.report.Bugs.BugStats = append(ts.report.Bugs.BugStats, stat)
			} else {
				values = append(values, "")
			}
		}
		values = append(values, fmt.Sprintf("%d", sum))
		ts.report.Bugs.BugCounts.Values = append(ts.report.Bugs.BugCounts.Values, values)
	}
}

func (ts *TicketStats) features() {
	features := FilterByType(ts.issues, "New Feature")
	openFeatures := OpenTickets(features)
	OrderByDue(openFeatures)

	for _, feature := range openFeatures {
		rf := feature.ToReportIssue(ts.jiraBase)
		if len(rf.Parents) == 0 {
			ts.report.Features = append(ts.report.Features, rf)
		}
	}
}

func (ts *TicketStats) improvements() {
	improvements := FilterByType(ts.issues, "Improvement")
	openImprovements := OpenTickets(improvements)
	OrderByDue(openImprovements)

	for _, improvement := range openImprovements {
		ri := improvement.ToReportIssue(ts.jiraBase)
		if len(ri.Parents) == 0 {
			ts.report.Improvements = append(ts.report.Improvements, ri)
		}
	}
}

func (ts *TicketStats) other() {
	others := Filter(ts.issues, func(issue *Issue) bool {
		return !(issue.Type == "Bug" || issue.Type == "New Feature" || issue.Type == "Improvement")
	})

	ts.report.Other.Count = len(OpenTickets(others))

	for _, t := range Types(others) {
		issues := FilterByType(others, t)
		count := len(OpenTickets(issues))

		statWeek := OtherTypeStats{
			Count: count,
			Type:  t,
		}
		statWeek.Report.Created = len(CreatedLastWeek(issues))
		statWeek.Report.Resolved = len(ClosedLastWeek(issues))
		statWeek.Report.Diff = statWeek.Report.Created - statWeek.Report.Resolved

		ts.report.Other.Week = append(ts.report.Other.Week, statWeek)

		statMonth := OtherTypeStats{
			Count: count,
			Type:  t,
		}
		statMonth.Report.Created = len(CreatedLastMonth(issues))
		statMonth.Report.Resolved = len(ClosedLastMonth(issues))
		statMonth.Report.Diff = statWeek.Report.Created - statWeek.Report.Resolved

		ts.report.Other.Month = append(ts.report.Other.Month, statMonth)
	}
}

func (ts *TicketStats) resources() {
	ranges := []string{"Last week", "Last month", "Last quarter", "Last year"}
	hours := calcHours(ts.issues)
	fte := calcFTE(hours)

	for i, r := range ranges {
		ts.report.Resources.Spend = append(ts.report.Resources.Spend, ResourceSpend{
			TimeRange: r,
			Effort:    formatWork(hours[i]),
			FTE:       fmt.Sprintf("%.2f", fte[i]),
		})
	}

	groups := make([]ResourceGroup, 4)
	for _, g := range groups {
		g.Type = "Type"
	}
	types := Types(ts.issues)
	sort.Slice(types, func(i, j int) bool {
		return strings.Compare(types[i], types[j]) < 0
	})
	for _, t := range types {
		issuesByType := FilterByType(ts.issues, t)

		ghours := calcHours(issuesByType)
		gfte := calcFTE(ghours)

		for i, g := range groups {
			percent := int((ghours[i] / hours[i]) * 100.0)
			if percent < 3 {
				continue
			}

			g.Details = append(g.Details, ResourceDetails{
				Type:    t,
				Work:    formatWork(ghours[i]),
				FTE:     fmt.Sprintf("%.2f", gfte[i]),
				Percent: percent,
			})
			groups[i] = g
		}
	}
	ts.report.Resources.Usage = append(ts.report.Resources.Usage, groups)

	groups = make([]ResourceGroup, 4)
	for _, g := range groups {
		g.Type = "Label"
	}
	labels := Labels(ts.issues)
	sort.Slice(labels, func(i, j int) bool {
		return strings.Compare(labels[i], labels[j]) < 0
	})
	for _, l := range labels {
		issuesByType := FilterByLabel(ts.issues, l)

		ghours := calcHours(issuesByType)
		gfte := calcFTE(ghours)

		for i, g := range groups {
			percent := int((ghours[i] / hours[i]) * 100.0)
			if percent < 5 {
				continue
			}
			g.Details = append(g.Details, ResourceDetails{
				Type:    l,
				Work:    formatWork(ghours[i]),
				FTE:     fmt.Sprintf("%.2f", gfte[i]),
				Percent: percent,
			})
			groups[i] = g
		}
	}
	ts.report.Resources.Usage = append(ts.report.Resources.Usage, groups)

	averageQuarter := NewResourceAverage()
	averageQuarter.TimeRange = "Last quarter"
	averageYear := NewResourceAverage()
	averageYear.TimeRange = "Last year"
	for issueType, times := range ResultionTimesByType(ClosedLastYear(ts.issues)) {
		averageQuarter.Details = append(averageQuarter.Details, ResourceAverageDetails{
			Type:   issueType,
			Count:  times.Quarter.Count,
			Median: formatWork(times.Quarter.Median),
			Mean:   formatWork(times.Quarter.Mean),
		})

		averageYear.Details = append(averageYear.Details, ResourceAverageDetails{
			Type:   issueType,
			Count:  times.Year.Count,
			Median: formatWork(times.Year.Median),
			Mean:   formatWork(times.Year.Mean),
		})
	}
	ts.report.Resources.Average = append(ts.report.Resources.Average, averageQuarter, averageYear)
}

func calcHours(issues []*Issue) []Work {
	hours := make([]Work, 0)
	hours = append(hours, WorkAfter(issues, time.Now().AddDate(0, 0, -7)))
	hours = append(hours, WorkAfter(issues, time.Now().AddDate(0, -1, 0)))
	hours = append(hours, WorkAfter(issues, time.Now().AddDate(0, -3, 0)))
	hours = append(hours, WorkAfter(issues, time.Now().AddDate(-1, 0, 0)))
	return hours
}

func calcFTE(hours []Work) []float64 {
	fte := make([]float64, 0)
	fte = append(fte, float64(hours[0]/40.0))
	fte = append(fte, float64(hours[1]/(40.0*4.25)))
	fte = append(fte, float64(hours[2]/(40.0*4.25*3)))
	fte = append(fte, float64(hours[3]/(40.0*52)))
	return fte
}
