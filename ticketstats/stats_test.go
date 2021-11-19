package ticketstats

import (
	"log"
	"testing"
	"time"
)

func TestOldBugs(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Type = "Bug"
	issue.Created = time.Now().AddDate(0, 0, -5)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Type = "Bug"
	issue.Created = time.Now().AddDate(0, 0, -25)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Type = "Bug"
	issue.Created = time.Now().AddDate(0, -1, -1)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "D"
	issue.Type = "Task"
	issue.Created = time.Now().AddDate(0, -1, -1)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "A"
	issue.Type = "Bug"
	issue.Resolved = time.Now().AddDate(0, 0, -1)
	issue.Created = time.Now().AddDate(0, -1, -5)
	issues = append(issues, issue)

	old := OldBugs(issues, DefaultConfig())

	if len(old) != 1 {
		log.Println("TEST: wrong count of old bugs")
		t.Fail()
	}

	if old[0].Key != "C" {
		log.Println("TEST: wrong issue")
		t.Fail()
	}
}

func TestResolutionTime(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.TimeSpend = Work(20.0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(30.0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(40.0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(300.0)
	issues = append(issues, issue)

	issue = NewIssue()
	issues = append(issues, issue)

	stats := ResolutionTime(issues)

	if stats.Count != 4 {
		t.Fail()
	}
	if stats.Mean != 97.5 {
		t.Fail()
	}
	if stats.Median != 35.0 {
		t.Fail()
	}
}

func TestResultionTimesByType(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.TimeSpend = Work(10.0)
	issue.Type = "Bug"
	issue.Resolved = time.Now().AddDate(0, -9, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(20.0)
	issue.Type = "Feature"
	issue.Resolved = time.Now().AddDate(0, -2, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(30.0)
	issue.Type = "Improvement"
	issue.Resolved = time.Now().AddDate(0, 0, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.TimeSpend = Work(40.0)
	issue.Type = "Task"
	issue.Resolved = time.Now().AddDate(0, 0, -20)
	issues = append(issues, issue)

	stats := ResultionTimesByType(issues)

	bug := stats["Bug"]
	if bug.Year.Count != 1 ||
		bug.Year.Mean != Work(10.0) ||
		bug.Year.Median != Work(10.0) {
		t.Fail()
	}
	if bug.Quarter.Count != 0 ||
		bug.Quarter.Mean != Work(0.0) ||
		bug.Quarter.Median != Work(0.0) {
		t.Fail()
	}
	if bug.Month.Count != 0 ||
		bug.Month.Mean != Work(0.0) ||
		bug.Month.Median != Work(0.0) {
		t.Fail()
	}
	if bug.Week.Count != 0 ||
		bug.Week.Mean != Work(0.0) ||
		bug.Week.Median != Work(0.0) {
		t.Fail()
	}

	feature := stats["Feature"]
	if feature.Year.Count != 1 ||
		feature.Year.Mean != Work(20.0) ||
		feature.Year.Median != Work(20.0) {
		t.Fail()
	}
	if feature.Quarter.Count != 1 ||
		feature.Quarter.Mean != Work(20.0) ||
		feature.Quarter.Median != Work(20.0) {
		t.Fail()
	}
	if feature.Month.Count != 0 ||
		feature.Month.Mean != Work(0.0) ||
		feature.Month.Median != Work(0.0) {
		t.Fail()
	}
	if feature.Week.Count != 0 ||
		feature.Week.Mean != Work(0.0) ||
		feature.Week.Median != Work(0.0) {
		t.Fail()
	}

	improvement := stats["Improvement"]
	if improvement.Year.Count != 1 ||
		improvement.Year.Mean != Work(30.0) ||
		improvement.Year.Median != Work(30.0) {
		t.Fail()
	}
	if improvement.Quarter.Count != 1 ||
		improvement.Quarter.Mean != Work(30.0) ||
		improvement.Quarter.Median != Work(30.0) {
		t.Fail()
	}
	if improvement.Month.Count != 1 ||
		improvement.Month.Mean != Work(30.0) ||
		improvement.Month.Median != Work(30.0) {
		t.Fail()
	}
	if improvement.Week.Count != 1 ||
		improvement.Week.Mean != Work(30.0) ||
		improvement.Week.Median != Work(30.0) {
		t.Fail()
	}

	task := stats["Task"]
	if task.Year.Count != 1 ||
		task.Year.Mean != Work(40.0) ||
		task.Year.Median != Work(40.0) {
		t.Fail()
	}
	if task.Quarter.Count != 1 ||
		task.Quarter.Mean != Work(40.0) ||
		task.Quarter.Median != Work(40.0) {
		t.Fail()
	}
	if task.Month.Count != 1 ||
		task.Month.Mean != Work(40.0) ||
		task.Month.Median != Work(40.0) {
		t.Fail()
	}
	if task.Week.Count != 0 ||
		task.Week.Mean != Work(0.0) ||
		task.Week.Median != Work(0.0) {
		t.Fail()
	}
}

func TestWorkAfter(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    10.0,
		Date:     time.Now().AddDate(0, 0, -2),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    20.0,
		Date:     time.Now().AddDate(0, 0, -4),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    30.0,
		Date:     time.Now().AddDate(0, 0, -10),
		Activity: "123456",
	})
	issues = append(issues, issue)

	issue = NewIssue()
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    40.0,
		Date:     time.Now().AddDate(0, 0, -3),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    50.0,
		Date:     time.Now().AddDate(0, 0, -11),
		Activity: "123456",
	})
	issues = append(issues, issue)

	work := WorkAfter(issues, time.Now().AddDate(0, 0, -7))

	if work != Work(70.0) {
		t.Fail()
	}
}
