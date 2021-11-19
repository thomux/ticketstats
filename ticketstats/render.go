package ticketstats

import (
	"bufio"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"time"
)

//go:embed report.tmpl
var reportTemplate string

// Report groups all data needed to render the HTML report.
type Report struct {
	Component    string
	Date         string
	OldBugs      []ReportIssue
	Bugs         ReportBugs
	Features     []ReportIssue
	Improvements []ReportIssue
	OtherCount   int
	Other        OtherReport
	Resources    ResourceReport
	HasWarnings  bool
	Warnings     Warnings
}

// NewReport initializes a new Report.
func NewReport() Report {
	var report Report

	report.OldBugs = make([]ReportIssue, 0)
	report.Bugs = NewReportBugs()
	report.Features = make([]ReportIssue, 0)
	report.Improvements = make([]ReportIssue, 0)
	report.Other = NewOtherReport()
	report.Resources = NewResourceReport()
	report.HasWarnings = false
	report.Warnings = NewWarnings()

	return report
}

// Warnings groups all sanitize warnings.
type Warnings struct {
	Count          int
	NoActivity     []ReportIssue
	InvalidBooking []InvalidBooking
}

func (sr SanitizeResult) ToWarnings(jiraBaseUrl string) Warnings {
	warnings := NewWarnings()
	warnings.Count = len(sr.NoActivity) + len(sr.InvalidWorkLogs)
	for _, na := range sr.NoActivity {
		warnings.NoActivity = append(warnings.NoActivity, na.ToReportIssue(jiraBaseUrl))
	}
	for _, il := range sr.InvalidWorkLogs {
		ib := NewInvalidBooking()
		ib.Issue = il.Issue.ToReportIssue(jiraBaseUrl)
		for _, wl := range il.Logs {
			ib.Logs = append(ib.Logs, InvalidLog{
				Activity: wl.Activity,
				Date:     wl.Date.Format("2006-01-02"),
				Effort:   formatWork(wl.Hours),
			})
		}
		warnings.InvalidBooking = append(warnings.InvalidBooking, ib)
	}
	return warnings
}

// NewWarnings initializes a new Warnings object.
func NewWarnings() Warnings {
	var warnings Warnings

	warnings.NoActivity = make([]ReportIssue, 0)
	warnings.InvalidBooking = make([]InvalidBooking, 0)

	return warnings
}

// InvalidBooking represents an invalid time recording.
type InvalidBooking struct {
	Issue ReportIssue
	Logs  []InvalidLog
}

// NewInvalidBooking initializes a new Warnings InvalidBooking.
func NewInvalidBooking() InvalidBooking {
	var booking InvalidBooking

	booking.Logs = make([]InvalidLog, 0)

	return booking
}

// InvalidLog groups the data of an (invalid) time log.
type InvalidLog struct {
	Activity string
	Date     string
	Effort   string
}

// ResourceReport groups the data on spend working hours.
type ResourceReport struct {
	Spend   []ResourceSpend
	Usage   [][]ResourceGroup
	Average []ResourceAverage
}

// NewResourceReport initializes a new ResourceReport.
func NewResourceReport() ResourceReport {
	var report ResourceReport

	report.Spend = make([]ResourceSpend, 0)
	report.Usage = make([][]ResourceGroup, 0)
	report.Average = make([]ResourceAverage, 0)

	return report
}

// ResourceAverage groups the information how much time was
// spend on a ticket type in average.
type ResourceAverage struct {
	TimeRange string
	Details   []ResourceAverageDetails
}

// NewResourceAverage initializes a new ResourceAverage object.
func NewResourceAverage() ResourceAverage {
	var avg ResourceAverage

	avg.Details = make([]ResourceAverageDetails, 0)

	return avg
}

// ResourceSpend groups the effort spend on a time range.
type ResourceSpend struct {
	TimeRange string
	Effort    string
	FTE       string
}

// ResourceGroup groups the ResourceDetails with a type name.
type ResourceGroup struct {
	Type    string
	Details []ResourceDetails
}

// NewResourceGroup initializes a new ResourceGroup.
func NewResourceGroup() ResourceGroup {
	var rg ResourceGroup

	rg.Details = make([]ResourceDetails, 0)

	return rg
}

// ResourceDetails groups effort spend a type.
type ResourceDetails struct {
	Type    string
	Work    string
	FTE     string
	Percent int
}

// ResourceAverageDetails groups the data on average resource
// usage for a type.
type ResourceAverageDetails struct {
	Type   string
	Median string
	Mean   string
	Count  int
}

// Link groups the data for web link.
type Link struct {
	Name string
	Url  string
}

// ReportIssue groups all data about a Jira issue needed for
// rendering the report.
type ReportIssue struct {
	JiraUrl     string
	Key         string
	Summary     string
	Activity    string
	Priority    string
	HasDue      bool
	Due         string
	Created     string
	Age         int
	Labels      []string
	Creator     string
	Assignee    string
	Status      string
	FixVersions []string
	Estimate    string
	HasEstimate bool
	TimeSpend   string
	HasTime     bool
	Progress    int
	AtRisk      bool
	FTE         string
	HasChilds   bool
	Overtime    bool
	Childs      []ReportIssue
	Parents     []Link
}

func (issue *Issue) ToReportIssue(jiraBaseUrl string) ReportIssue {
	var rissue ReportIssue
	var noDate time.Time

	if jiraBaseUrl != "" {
		rissue.JiraUrl = jiraBaseUrl + issue.Key
	}
	rissue.Key = issue.Key
	rissue.Summary = issue.Summary
	rissue.Activity = issue.CustomActivity
	rissue.Priority = issue.Priority
	if issue.Due == noDate {
		rissue.HasDue = false
	} else {
		rissue.HasDue = true
		rissue.Due = issue.Due.Format("2006-01-02")
		if issue.OriginalEstimate > 0.1 {
			rissue.FTE = covertToFTE(issue.Due, issue.OriginalEstimate-issue.TimeSpend)
			rissue.HasEstimate = true
		}
	}
	if issue.Created != noDate {
		rissue.Created = issue.Created.Format("2006-01-02")
		rissue.Age = convertToAge(issue.Created)
	}
	rissue.Labels = issue.Labels
	rissue.Creator = issue.Creator
	rissue.Assignee = issue.Assignee
	rissue.Status = issue.Status
	rissue.FixVersions = issue.FixVersions
	if issue.OriginalEstimate > 0.001 {
		rissue.Estimate = formatWork(issue.OriginalEstimate)
	}
	if issue.TimeSpend > 0.1 {
		rissue.TimeSpend = formatWork(issue.TimeSpend)
	}
	if issue.OriginalEstimate > 0.1 && issue.TimeSpend > 0.1 {
		rissue.HasTime = true
		rissue.Progress = int((issue.TimeSpend / issue.OriginalEstimate) * 100.0)
		if issue.TimeSpend > issue.OriginalEstimate {
			rissue.Overtime = true
		}
	}
	if len(issue.Childs) > 0 {
		parent := Link{
			Url:  jiraBaseUrl + issue.Key,
			Name: issue.Key,
		}
		rissue.Childs = flattenTree(issue, parent, jiraBaseUrl)
		rissue.HasChilds = (len(rissue.Childs) > 0)
	}

	return rissue
}

func flattenTree(issue *Issue, parent Link, jiraBaseUrl string) []ReportIssue {
	childs := make([]ReportIssue, 0)

	for _, child := range issue.Childs {
		rissue := child.ToReportIssue(jiraBaseUrl)
		rissue.Parents = append(rissue.Parents, parent)
		if rissue.Status != "Closed" {
			childs = append(childs, rissue)
		}
		for _, rc := range rissue.Childs {
			if rc.Status != "Closed" {
				childs = append(childs, rc)
			}
		}
	}

	return childs
}

func covertToFTE(due time.Time, remainingEffort Work) string {
	neededDays := float64(remainingEffort) / 8.0
	remainingTime := time.Until(due)
	remainingWeeks := (remainingTime.Hours() / 24.0) / 7.0
	remainingDays := remainingWeeks * 5
	fte := neededDays / float64(remainingDays)
	return fmt.Sprintf("%.2f", fte)
}

func convertToAge(date time.Time) int {
	duration := time.Since(date)
	return int(duration.Hours()) / 24
}

// NewReportIssue initializes a new ReportIssue.
func NewReportIssue() ReportIssue {
	var issue ReportIssue

	issue.HasDue = false
	issue.HasTime = false
	issue.Labels = make([]string, 0)
	issue.FixVersions = make([]string, 0)
	issue.HasChilds = false
	issue.Childs = make([]ReportIssue, 0)
	issue.Parents = make([]Link, 0)
	issue.AtRisk = false

	return issue
}

// ReportBugs groups the data for the bug report section.
type ReportBugs struct {
	Count     int
	Week      ReportCount
	Month     ReportCount
	BugStats  []ReportBugStats
	BugCounts BugCounts
}

// NewReportBugs initializes a new ReportBugs object.
func NewReportBugs() ReportBugs {
	var report ReportBugs

	report.BugStats = make([]ReportBugStats, 0)

	return report
}

// BugCounts groups the bug numbers by version
type BugCounts struct {
	Versions []string
	Values   [][]string
}

// BugCounts initializes a new BugCounts object.
func NewBugCounts() BugCounts {
	var bugCounts BugCounts
	bugCounts.Versions = make([]string, 0)
	bugCounts.Values = make([][]string, 0)
	return bugCounts
}

// ReportCount groups the count changes for a issues type.
type ReportCount struct {
	Created  int
	Resolved int
	Diff     int
}

// ReportBugStats groups the bug statistics for a fix version.
type ReportBugStats struct {
	Version  string
	Security string
	Count    int
	Bugs     []ReportIssue
}

// NewReportBugStats initializes a new ReportBugStats object.
func NewReportBugStats() ReportBugStats {
	var report ReportBugStats

	report.Bugs = make([]ReportIssue, 0)

	return report
}

// OtherReport groups the data for the "other issues" section.
type OtherReport struct {
	Count int
	Week  []OtherTypeStats
	Month []OtherTypeStats
}

// NewOtherReport initializes a new OtherReport object.
func NewOtherReport() OtherReport {
	var report OtherReport

	report.Week = make([]OtherTypeStats, 0)
	report.Month = make([]OtherTypeStats, 0)

	return report
}

// OtherTypeStats groups the data for "other issue" types.
type OtherTypeStats struct {
	Type   string
	Count  int
	Report ReportCount
}

// Report.Render renders an HTMl report.
func (report Report) Render() {
	path := "./report_" + report.Component + ".html"

	_ = os.Remove(path)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)

	t := template.New("report")
	sec := false
	t.Funcs(template.FuncMap{"second": func() bool {
		sec = !sec
		return sec
	}})
	t, err = t.Parse(reportTemplate)
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(w, "report", report)
	if err != nil {
		panic(err)
	}

	w.Flush()
	f.Close()
}
