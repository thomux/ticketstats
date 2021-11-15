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

	// read issues form csv
	issues := Parse(path)

	if project != "" {
		issues = FilterByProject(issues, project)
	}
	if component != "" {
		issues = FilterByComponent(issues, component)
		splitByComponent = false
	}

	ts := TicketStats{
		jiraBase: jiraBase,
		issues:   issues,
		report:   NewReport(),
	}
	ts.report.Component = component
	ts.report.Date = time.Now().Format("2006-01-02")
	ts.ignoreOld = true
	ts.generateReport()

	if splitByComponent {
		log.Println("ERROR: split by component not implemented")
		// TODO: generte reports for components
	}
}

func (ts *TicketStats) generateReport() {
	// Reduce to active tickets
	ts.active = ActiveTickets(ts.issues)
	log.Println("INFO:", len(ts.active), "active tickets.")

	ts.sanitize()
	ts.oldBugs()
	ts.bugs()

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
	OrderByCreated(oldBugs)
	for _, bug := range oldBugs {
		ts.report.OldBugs = append(ts.report.OldBugs, bug.ToReportIssue(ts.jiraBase))
	}
	log.Println("INFO:", len(oldBugs), "old bug tickets.")
}

func (ts *TicketStats) bugs() {
	bugs := FilterByType(ts.issues, "Bug")
	openBugs := OpenTickets(bugs)

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

				for _, b := range FilterByPriority(bs, "Critical") {
					stat.Critical = append(stat.Critical, b.ToReportIssue(ts.jiraBase))
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
