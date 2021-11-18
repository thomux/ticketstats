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
	issueA := NewIssue()
	issueA.Key = "A"

	issueB := NewIssue()
	issueB.Key = "B"

	issueC := NewIssue()
	issueC.Key = "C"
	issueC.Childs = append(issueC.Childs, issueA, issueB)

	issueD := NewIssue()
	issueD.Key = "D"

	issueE := NewIssue()
	issueE.Key = "E"
	issueE.Childs = append(issueE.Childs, issueC)

	issueF := NewIssue()
	issueF.Key = "F"
	issueF.Childs = append(issueF.Childs, issueD, issueE)

	linkParentsRecursive(issueF)

	if len(issueE.Parents) != 1 {
		log.Println("TEST: issueE parents len wrong")
		t.Fail()
	}
	if issueE.Parents[0].Key != "F" {
		log.Println("TEST: issueE parent wrong")
		t.Fail()
	}

	if len(issueD.Parents) != 1 {
		log.Println("TEST: issueD parents len wrong")
		t.Fail()
	}
	if issueD.Parents[0].Key != "F" {
		log.Println("TEST: issueD parent wrong")
		t.Fail()
	}

	if len(issueC.Parents) != 1 {
		log.Println("TEST: issueD parents len wrong")
		t.Fail()
	}
	if issueC.Parents[0].Key != "E" {
		log.Println("TEST: issueD parent wrong")
		t.Fail()
	}

	if len(issueB.Parents) != 1 {
		log.Println("TEST: issueC parents len wrong")
		t.Fail()
	}
	if issueB.Parents[0].Key != "C" {
		log.Println("TEST: issueC parent wrong")
		t.Fail()
	}

	if len(issueA.Parents) != 1 {
		log.Println("TEST: issueA parents len wrong")
		t.Fail()
	}
	if issueA.Parents[0].Key != "C" {
		log.Println("TEST: issueA parent wrong")
		t.Fail()
	}
}

func TestRemoveDuplicateChildsRecursive(t *testing.T) {
	t.Fail()
}
