package ticketstats

import (
	"testing"
	"time"
)

func TestLoadTemplate(t *testing.T) {
	temp := loadTemplate(DefaultConfig())
	if temp == nil {
		t.Fail()
	}
	if temp.Name() != "report" {
		t.Fail()
	}
}

func TestToWarnings(t *testing.T) {
	noActiviy := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Summary = "B"
	noActiviy = append(noActiviy, issue)

	logs := make([]WorkLog, 0)
	logs = append(logs, WorkLog{
		Hours:    Work(123.0),
		Date:     time.Now(),
		Activity: "123456",
	})

	invalidWork := make([]InvalidWorkLog, 0)
	invalidWork = append(invalidWork, InvalidWorkLog{
		Issue: issue,
		Logs:  logs,
	})

	sr := SanitizeResult{
		NoActivity:      noActiviy,
		InvalidWorkLogs: invalidWork,
	}

	w := sr.ToWarnings("https://test.url/", DefaultConfig())

	if w.Count != 2 {
		t.Fail()
	}

	if len(w.NoActivity) != 1 {
		t.Fail()
	}

	if len(w.InvalidBooking) != 1 {
		t.Fail()
	}
}

func TestToReportIssue(t *testing.T) {
	ccissue := NewIssue()
	ccissue.Key = "C"

	cissue := NewIssue()
	cissue.Key = "B"
	cissue.Childs = append(cissue.Childs, ccissue)

	cissue2 := NewIssue()
	cissue2.Key = "D"

	issue := NewIssue()
	issue.Key = "A"
	issue.Summary = "summary"
	issue.CustomActivity = "123456"
	issue.Priority = "prio"
	issue.Due = time.Now()
	issue.OriginalEstimate = Work(8.0)
	issue.TimeSpend = Work(4.0)
	issue.Created = time.Now().AddDate(0, 0, -7)
	issue.Labels = append(issue.Labels, "label")
	issue.Creator = "aCreator"
	issue.Assignee = "aAssignee"
	issue.Status = "Open"
	issue.FixVersions = append(issue.FixVersions, "1.2.3")
	issue.Childs = append(issue.Childs, cissue, cissue2)

	ri := issue.ToReportIssue("https://test.url/", DefaultConfig())

	if ri.Key != "A" {
		t.Fail()
	}
	if ri.Summary != "summary" {
		t.Fail()
	}
	if ri.Activity != "123456" {
		t.Fail()
	}
	if ri.Priority != "prio" {
		t.Fail()
	}
	if ri.Due == "" {
		t.Fail()
	}
	if !ri.HasDue {
		t.Fail()
	}
	if ri.FTE == "" {
		t.Fail()
	}
	if !ri.HasEstimate {
		t.Fail()
	}
	if ri.Created == "" {
		t.Fail()
	}
	if ri.Age != 7 {
		t.Fail()
	}
	if len(ri.Labels) != 1 {
		t.Fail()
	}
	if ri.Labels[0] != "label" {
		t.Fail()
	}
	if ri.Creator != "aCreator" {
		t.Fail()
	}
	if ri.Assignee != "aAssignee" {
		t.Fail()
	}
	if ri.Status != "Open" {
		t.Fail()
	}
	if len(ri.FixVersions) != 1 {
		t.Fail()
	}
	if ri.FixVersions[0] != "1.2.3" {
		t.Fail()
	}
	if ri.Estimate == "" {
		t.Fail()
	}
	if ri.TimeSpend == "" {
		t.Fail()
	}
	if !ri.HasTime {
		t.Fail()
	}
	if ri.Progress != 50 {
		t.Fail()
	}
	if ri.Overtime {
		t.Fail()
	}
	if !ri.HasChilds {
		t.Fail()
	}
	if len(ri.Childs) != 3 {
		t.Fail()
	}
	if ri.Childs[0].Key != "B" {
		t.Fail()
	}
	if ri.Childs[1].Key != "C" {
		t.Fail()
	}
	if ri.Childs[2].Key != "D" {
		t.Fail()
	}
	if ri.JiraUrl != "https://test.url/A" {
		t.Fail()
	}
}

func TestFlattenTree(t *testing.T) {
	ccissue := NewIssue()
	ccissue.Key = "C"

	cissue := NewIssue()
	cissue.Key = "B"
	cissue.Childs = append(cissue.Childs, ccissue)

	cissue2 := NewIssue()
	cissue2.Key = "D"

	issue := NewIssue()
	issue.Key = "A"
	issue.Summary = "summary"
	issue.CustomActivity = "123456"
	issue.Priority = "prio"
	issue.Due = time.Now()
	issue.OriginalEstimate = Work(8.0)
	issue.TimeSpend = Work(4.0)
	issue.Created = time.Now().AddDate(0, 0, -7)
	issue.Labels = append(issue.Labels, "label")
	issue.Creator = "aCreator"
	issue.Assignee = "aAssignee"
	issue.Status = "Open"
	issue.FixVersions = append(issue.FixVersions, "1.2.3")
	issue.Childs = append(issue.Childs, cissue, cissue2)

	parent := Link{
		Name: "Parent",
		Url:  "https://test.url/",
	}

	childs := flattenTree(issue, parent, "https://jira.url/", DefaultConfig())

	if len(childs) != 3 {
		t.Fail()
	}
	if childs[0].Key != "B" ||
		childs[1].Key != "C" ||
		childs[2].Key != "D" {
		t.Fail()
	}
}

func TestCovertToFTE(t *testing.T) {
	fte := covertToFTE(time.Now().AddDate(0, 0, 14), Work(40.0))
	if fte != "0.50" {
		t.Fail()
	}
}

func TestConvertToAge(t *testing.T) {
	days := convertToAge(time.Now().AddDate(0, 0, -20))
	if days != 20 {
		t.Fail()
	}
}
