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
	t.Fail()
}

func TestResultionTimesByType(t *testing.T) {
	t.Fail()
}

func TestContainedTypes(t *testing.T) {
	t.Fail()
}

func TestWorkAfter(t *testing.T) {
	t.Fail()
}
