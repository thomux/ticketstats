// ticketstats generates statistics and further report for a set of csv exported
// Jira tickets.
package ticketstats

import (
	"log"
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
