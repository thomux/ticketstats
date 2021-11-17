package ticketstats

import (
	"log"
	"testing"
	"time"
)

func TestRemoveDuplicates(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.LinkParents = append(issue.LinkParents, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.LinkParents = append(issue.LinkParents, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Id = "C"
	issues = append(issues, issue)

	issues = removeDuplicates(issues, make(map[string]*Issue))

	if len(issues) != 3 {
		log.Println("TEST: wrong issue count")
		t.Fail()
	}

	keys := ""
	for _, issue := range issues {
		keys += issue.Key
	}
	if keys != "ABC" {
		log.Println("TEST: wrong issues")
		t.Fail()
	}
}

func TestClusterIssuesParents(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.LinkParents = append(issue.LinkParents, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[1].Childs) != 0 {
		log.Println("TEST: child issue")
		t.Fail()
	}

	if len(issues[0].Childs) != 1 {
		log.Println("TEST: parent issue")
		t.Fail()
	}
}

func TestClusterIssuesParts(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.LinkParts = append(issue.LinkParts, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[1].Childs) != 1 {
		log.Println("TEST: parent issue")
		t.Fail()
	}

	if len(issues[0].Childs) != 0 {
		log.Println("TEST: child issue")
		t.Fail()
	}
}

func TestClusterIssuesSplits(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.LinkIssueSplits = append(issue.LinkIssueSplits, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[0].Childs) != 1 {
		log.Println("TEST: parent issue")
		t.Fail()
	}

	if len(issues[1].Childs) != 0 {
		log.Println("TEST: child issue")
		t.Fail()
	}
}

func TestClusterIssuesClones(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.Created = time.Now().AddDate(0, 0, -7)
	issue.LinkCloners = append(issue.LinkCloners, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issue.Created = time.Now().AddDate(0, 0, -5)
	issue.LinkCloners = append(issue.LinkCloners, "A")
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[0].Childs) != 1 {
		log.Println("TEST: older issue should be parent")
		t.Fail()
	}

	if len(issues[1].Childs) != 0 {
		log.Println("TEST: younger issue should be child")
		t.Fail()
	}
}

func TestClusterIssuesDuplicatesNotClosed(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.Created = time.Now().AddDate(0, 0, -7)
	issue.LinkDuplicates = append(issue.LinkDuplicates, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issue.Created = time.Now().AddDate(0, 0, -5)
	issue.LinkDuplicates = append(issue.LinkDuplicates, "A")
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[0].Childs) != 1 {
		log.Println("TEST: older issue should be parent")
		t.Fail()
	}

	if len(issues[1].Childs) != 0 {
		log.Println("TEST: younger issue should be child")
		t.Fail()
	}
}

func TestClusterIssuesDuplicatesClosed(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Id = "A"
	issue.Created = time.Now().AddDate(0, 0, -7)
	issue.LinkDuplicates = append(issue.LinkDuplicates, "B")
	issue.Status = "Closed"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Id = "B"
	issue.Created = time.Now().AddDate(0, 0, -5)
	issue.LinkDuplicates = append(issue.LinkDuplicates, "A")
	issues = append(issues, issue)

	ClusterIssues(issues)

	if len(issues[0].Childs) != 0 {
		log.Println("TEST: closed issue should be child")
		t.Fail()
	}

	if len(issues[1].Childs) != 1 {
		log.Println("TEST: open issue should be parent")
		t.Fail()
	}
}

func TestLinkParentsRecursive(t *testing.T) {
	t.Fail()
}

func TestRemoveDuplicateChildsRecursive(t *testing.T) {
	t.Fail()
}
