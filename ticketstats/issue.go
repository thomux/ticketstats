package ticketstats

import (
	"fmt"
	"time"
)

// Work is the data type or logged work time.
// The unit is hours.
type Work float64

// WorkLog groups the infos from one logged work time.
type WorkLog struct {
	// worked hours
	Hours Work
	// date of time recording
	Date time.Time
	// activity (custom value)
	Activity string
}

// formatWork converts worked hours to a string.
// The representation splits the hours to weeks, days and hours
// for better readability.
func formatWork(w Work) string {
	if w < 8 {
		// short cut for less than one day
		return fmt.Sprintf("%.2fh", w)
	}

	days := int(w / 8)
	hours := float64(w) - float64(days*8)
	weeks := days / 5
	days = days % 5

	str := ""

	if weeks > 0 {
		str += fmt.Sprintf("%dw ", weeks)
	}
	if days > 0 {
		str += fmt.Sprintf("%dd ", days)
	}
	if hours > 0.1 {
		str += fmt.Sprintf("%.2fh", hours)
	}

	return str
}

// WorkLog.ToString converts a WorkLog to a string for printing.
func (workLog WorkLog) ToString() string {
	return fmt.Sprintf("%s: %s - %s\n",
		workLog.Activity,
		workLog.Date.Format("2006-01-02"),
		formatWork(workLog.Hours))
}

// Issue groups all needed Jira issue data.
type Issue struct {
	Summary              string
	Key                  string
	Id                   string
	Parent               string
	Type                 string
	Status               string
	Priority             string
	Assignee             string
	Creator              string
	Created              time.Time
	Updated              time.Time
	LastViewed           time.Time
	AffectsVersions      []string
	FixVersions          []string
	Components           []string
	LogWorks             []WorkLog
	OriginalEstimate     Work
	RemainingEstimate    Work
	TimeSpend            Work
	SumOriginalEstimate  Work
	SumRemainingEstimate Work
	SumTimeSpend         Work
	SecurityLevel        string
	Labels               []string
	Resolution           string
	Resolved             time.Time
	Due                  time.Time
	LinkBlocks           []string
	LinkCauses           []string
	LinkCloners          []string
	LinkDependencies     []string
	LinkDuplicates       []string
	LinkIssueSplits      []string
	LinkParts            []string
	LinkRelates          []string
	LinkRelations        []string
	LinkTriggers         []string
	LinkLinkIssues       []string
	LinkParents          []string
	CustomExternalId     string
	CustomSupplierRef    string
	CustomVariant        string
	CustomActivity       string
	CustomCategory       string
	Childs               []*Issue
	Parents              []*Issue
}

// NewIssue creates a new issue.
func NewIssue() *Issue {
	issue := new(Issue)

	// init all lists to avoid nil issues
	issue.AffectsVersions = make([]string, 0)
	issue.FixVersions = make([]string, 0)
	issue.Components = make([]string, 0)
	issue.LogWorks = make([]WorkLog, 0)
	issue.Labels = make([]string, 0)
	issue.LinkBlocks = make([]string, 0)
	issue.LinkCauses = make([]string, 0)
	issue.LinkCloners = make([]string, 0)
	issue.LinkDependencies = make([]string, 0)
	issue.LinkDuplicates = make([]string, 0)
	issue.LinkIssueSplits = make([]string, 0)
	issue.LinkParts = make([]string, 0)
	issue.LinkRelates = make([]string, 0)
	issue.LinkRelations = make([]string, 0)
	issue.LinkTriggers = make([]string, 0)
	issue.LinkLinkIssues = make([]string, 0)
	issue.LinkParents = make([]string, 0)
	issue.Childs = make([]*Issue, 0)
	issue.Parents = make([]*Issue, 0)

	return issue
}

// Issue.IsResolved tests if the issue is resolved.
func (issue *Issue) IsResolved() bool {
	return issue.Resolved != time.Time{}
}

// Issue.ToString creates a string representation of the issue for console.
func (issue *Issue) ToString() string {
	str := ""

	var noDate time.Time

	str += fmt.Sprintf("Summary: %s\n", issue.Summary)
	str += fmt.Sprintf("Key: %s\n", issue.Key)
	str += fmt.Sprintf("Id: %s\n", issue.Id)
	if issue.Parent != "" {
		str += fmt.Sprintf("Parent: %s\n", issue.Parent)
	}
	str += fmt.Sprintf("Type: %s\n", issue.Type)
	str += fmt.Sprintf("Status: %s\n", issue.Status)
	str += fmt.Sprintf("Priority: %s\n", issue.Priority)
	str += fmt.Sprintf("Assignee: %s\n", issue.Assignee)
	str += fmt.Sprintf("Creator: %s\n", issue.Creator)
	str += fmt.Sprintf("Created: %s\n", issue.Created.Format("2006-01-02"))
	str += fmt.Sprintf("Updated: %s\n", issue.Updated.Format("2006-01-02"))
	str += fmt.Sprintf("Last viewed: %s\n", issue.LastViewed.Format("2006-01-02"))
	if len(issue.AffectsVersions) > 0 {
		str += fmt.Sprintf("Affects versions: %+v\n", issue.AffectsVersions)
	}
	if len(issue.FixVersions) > 0 {
		str += fmt.Sprintf("Fix versions: %+v\n", issue.FixVersions)
	}
	str += fmt.Sprintf("Components: %+v\n", issue.Components)
	if len(issue.LogWorks) > 0 {
		str += "Work Logs:\n"
		for _, l := range issue.LogWorks {
			str += "- " + l.ToString()
		}
	}
	if issue.OriginalEstimate > 0 {
		str += fmt.Sprintf("Original estimate: %s\n",
			formatWork(issue.OriginalEstimate))
	}
	if issue.RemainingEstimate > 0 {
		str += fmt.Sprintf("Remaining estimate: %s\n",
			formatWork(issue.RemainingEstimate))
	}
	if issue.TimeSpend > 0 {
		str += fmt.Sprintf("Time spend: %s\n", formatWork(issue.TimeSpend))
	}
	if issue.SumOriginalEstimate > 0 {
		str += fmt.Sprintf("Sum original estimate: %s\n",
			formatWork(issue.SumOriginalEstimate))
	}
	if issue.SumRemainingEstimate > 0 {
		str += fmt.Sprintf("Sum remaining estimate: %s\n",
			formatWork(issue.SumRemainingEstimate))
	}
	if issue.SumTimeSpend > 0 {
		str += fmt.Sprintf("Sum time spend: %s\n",
			formatWork(issue.SumTimeSpend))
	}
	if len(issue.Labels) > 0 {
		str += fmt.Sprintf("Labels: %+v\n", issue.Labels)
	}
	if issue.Resolution != "" {
		str += fmt.Sprintf("Resolution: %s\n", issue.Resolution)
	}
	if issue.Resolved != noDate {
		str += fmt.Sprintf("Resolved: %s\n",
			issue.Resolved.Format("2006-01-02"))
	}
	if issue.Due != noDate {
		str += fmt.Sprintf("Due: %s\n", issue.Due.Format("2006-01-02"))
	}
	if len(issue.LinkBlocks) > 0 {
		str += fmt.Sprintf("Links block: %+v\n", issue.LinkBlocks)
	}
	if len(issue.LinkCauses) > 0 {
		str += fmt.Sprintf("Links causes: %+v\n", issue.LinkCauses)
	}
	if len(issue.LinkCloners) > 0 {
		str += fmt.Sprintf("Links cloners: %+v\n", issue.LinkCloners)
	}
	if len(issue.LinkDependencies) > 0 {
		str += fmt.Sprintf("Links dependencies: %+v\n", issue.LinkDependencies)
	}
	if len(issue.LinkDuplicates) > 0 {
		str += fmt.Sprintf("Links duplicates: %+v\n", issue.LinkDuplicates)
	}
	if len(issue.LinkIssueSplits) > 0 {
		str += fmt.Sprintf("Links issues splits: %+v\n", issue.LinkIssueSplits)
	}
	if len(issue.LinkParts) > 0 {
		str += fmt.Sprintf("Links parts: %+v\n", issue.LinkParts)
	}
	if len(issue.LinkRelates) > 0 {
		str += fmt.Sprintf("Links relates: %+v\n", issue.LinkRelates)
	}
	if len(issue.LinkTriggers) > 0 {
		str += fmt.Sprintf("Links triggers: %+v\n", issue.LinkTriggers)
	}
	if len(issue.LinkLinkIssues) > 0 {
		str += fmt.Sprintf("Links links issue: %+v\n", issue.LinkLinkIssues)
	}
	if len(issue.LinkParents) > 0 {
		str += fmt.Sprintf("Links parent: %+v\n", issue.LinkParents)
	}
	if issue.CustomExternalId != "" {
		str += fmt.Sprintf("External ID: %s\n", issue.CustomExternalId)
	}
	if issue.CustomSupplierRef != "" {
		str += fmt.Sprintf("Supplier reference: %s\n", issue.CustomSupplierRef)
	}
	if issue.CustomVariant != "" {
		str += fmt.Sprintf("Variant: %s\n", issue.CustomVariant)
	}
	if issue.CustomActivity != "" {
		str += fmt.Sprintf("Activity: %s\n", issue.CustomActivity)
	}
	if issue.CustomCategory != "" {
		str += fmt.Sprintf("Category: %s\n", issue.CustomCategory)
	}
	if len(issue.Childs) > 0 {
		str += fmt.Sprintf("Childs: %+v\n", issue.Childs)
	}

	return str
}
