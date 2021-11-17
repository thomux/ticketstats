package ticketstats

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestActiveTickets(t *testing.T) {
	issues := make([]*Issue, 0)

	// issue which was updated recently
	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(0, 0, -5)
	issue.Updated = time.Now()
	issues = append(issues, issue)

	// issue which was not resolved
	issue = NewIssue()
	issue.Key = "B"
	// issue.Resolved not resolved
	issues = append(issues, issue)

	// issue which was resolved a while ago and not recently updated
	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, -2, 0)
	issue.Updated = time.Now().AddDate(0, -1, -10)
	issues = append(issues, issue)

	// call function to test
	issues = ActiveTickets(issues)

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if !strings.Contains(keys, "A") ||
		!strings.Contains(keys, "B") ||
		strings.Contains(keys, "C") {

		log.Println("TEST: issue count not expected")
		t.Fail()
	}
}

func TestFilterByProject(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "ICAS1-12345"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "ICAS1-12346"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "OTHER-12347"
	issues = append(issues, issue)

	// call function to test
	issues = FilterByProject(issues, "ICAS1")

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if !strings.Contains(keys, "ICAS1-12345") ||
		!strings.Contains(keys, "ICAS1-12346") ||
		strings.Contains(keys, "OTHER-12347") {

		log.Println("TEST: issue count not expected")
		t.Fail()
	}
}

func TestCreatedLastWeek(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(0, 0, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, 0, -5)
	issues = append(issues, issue)

	// call function to test
	issues = CreatedLastWeek(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestCreatedLastMonth(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(0, -2, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, 0, -20)
	issues = append(issues, issue)

	// call function to test
	issues = CreatedLastMonth(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestCreatedLastQuarter(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(0, -3, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(0, -5, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, -2, -5)
	issues = append(issues, issue)

	issues = CreatedLastQuarter(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestCreatedLastYear(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(-1, 0, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(-1, -5, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, -10, -5)
	issues = append(issues, issue)

	// call function to test
	issues = CreatedLastYear(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestClosedLastWeek(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Resolved = time.Now().AddDate(0, 0, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, 0, -5)
	issues = append(issues, issue)

	// call function to test
	issues = ClosedLastWeek(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestClosedLastMonth(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Resolved = time.Now().AddDate(0, -2, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, 0, -20)
	issues = append(issues, issue)

	// call function to test
	issues = ClosedLastMonth(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestClosedLastQuarter(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(0, -3, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Resolved = time.Now().AddDate(0, -5, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, -2, -5)
	issues = append(issues, issue)

	// call function to test
	issues = ClosedLastQuarter(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestClosedLastYear(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(-1, 0, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Resolved = time.Now().AddDate(-1, -5, -8)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, -10, -5)
	issues = append(issues, issue)

	// call function to test
	issues = ClosedLastYear(issues)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	if issues[0].Key != "C" {
		log.Println("TEST: issue key not expected")
		t.Fail()
	}
}

func TestOlderThanOneMonth(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(0, -1, -2)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(-1, 0, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, 0, -25)
	issues = append(issues, issue)

	// call function to test
	issues = OlderThanOneMonth(issues)

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if !strings.Contains(keys, "A") ||
		!strings.Contains(keys, "B") ||
		strings.Contains(keys, "C") {

		log.Println("TEST: issue count not expected")
		t.Fail()
	}
}

func TestFilterByType(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Type = "Bug"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Type = "Bug"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Type = "Task"
	issues = append(issues, issue)

	// call function to test
	issues = FilterByType(issues, "Bug")

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if !strings.Contains(keys, "A") ||
		!strings.Contains(keys, "B") ||
		strings.Contains(keys, "C") {

		log.Println("TEST: wrong issues")
		t.Fail()
	}
}

func TestFilterByComponent(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Components = append(issue.Components, "Linux")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Components = append(issue.Components, "Linux")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Components = append(issue.Components, "Hypervisor")
	issues = append(issues, issue)

	// call function to test
	issues = FilterByComponent(issues, "Linux")

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if !strings.Contains(keys, "A") ||
		!strings.Contains(keys, "B") ||
		strings.Contains(keys, "C") {

		log.Println("TEST: wrong issues")
		t.Fail()
	}
}

func TestOrderByCreated(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Created = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Created = time.Now().AddDate(0, -1, -1)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Created = time.Now().AddDate(0, 0, -20)
	issues = append(issues, issue)

	// call function to test
	OrderByCreated(issues)

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if keys != "BAC" {
		log.Println("TEST: wrong order")
		t.Fail()
	}
}

func TestAge(t *testing.T) {
	if Age(time.Now().AddDate(0, 0, -23)) != 23 {
		log.Println("TEST: wrong number of days")
		t.Fail()
	}
}

func TestIssuesByComponent(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Components = append(issue.Components, "Linux")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Components = append(issue.Components, "Linux")
	issue.Components = append(issue.Components, "Hypervisor")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Components = append(issue.Components, "Hypervisor")
	issues = append(issues, issue)

	// call function to test
	issuesByComponent := IssuesByComponent(issues)

	if len(issuesByComponent) != 2 {
		log.Println("TEST: component count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issuesByComponent["Linux"] {
		keys += i.Key
	}

	if !strings.Contains(keys, "A") ||
		!strings.Contains(keys, "B") ||
		strings.Contains(keys, "C") {

		log.Println("TEST: wrong issues")
		t.Fail()
	}

	keys = ""
	for _, i := range issuesByComponent["Hypervisor"] {
		keys += i.Key
	}

	if !strings.Contains(keys, "B") ||
		!strings.Contains(keys, "C") ||
		strings.Contains(keys, "A") {

		log.Println("TEST: wrong issues Hypervisor")
		t.Fail()
	}
}

func TestLastWeek(t *testing.T) {
	if !lastWeek(time.Now().AddDate(0, 0, -5)) {
		log.Println("TEST: date within last week")
		t.Fail()
	}

	if lastWeek(time.Now().AddDate(0, 0, -10)) {
		log.Println("TEST: date not within last week")
		t.Fail()
	}
}

func TestLastMonth(t *testing.T) {
	if !lastMonth(time.Now().AddDate(0, 0, -25)) {
		log.Println("TEST: date within last month")
		t.Fail()
	}

	if lastMonth(time.Now().AddDate(0, -1, -1)) {
		log.Println("TEST: date not within last month")
		t.Fail()
	}
}

func TestLastQuarter(t *testing.T) {
	if !lastQuarter(time.Now().AddDate(0, -2, -25)) {
		log.Println("TEST: date within last quarter")
		t.Fail()
	}

	if lastQuarter(time.Now().AddDate(0, -3, -1)) {
		log.Println("TEST: date not within last quarter")
		t.Fail()
	}
}

func TestLastYear(t *testing.T) {
	if !lastYear(time.Now().AddDate(0, -11, -25)) {
		log.Println("TEST: date within last year")
		t.Fail()
	}

	if lastYear(time.Now().AddDate(-1, 0, -1)) {
		log.Println("TEST: date not within last year")
		t.Fail()
	}
}

func TestOpenTickets(t *testing.T) {
	t.Fail()
}

func TestFixVersions(t *testing.T) {
	t.Fail()
}

func TestSecurityLevels(t *testing.T) {
	t.Fail()
}

func TestFilterByFixVersion(t *testing.T) {
	t.Fail()
}

func TestFilterBySecurityLevel(t *testing.T) {
	t.Fail()
}

func TestFilterByPriority(t *testing.T) {
	t.Fail()
}

func TestFilter(t *testing.T) {
	t.Fail()
}

func TestOrderByDue(t *testing.T) {
	t.Fail()
}

func TestOrderByPriority(t *testing.T) {
	t.Fail()
}

func TestOrderByStatus(t *testing.T) {
	t.Fail()
}

func TestTypes(t *testing.T) {
	t.Fail()
}
