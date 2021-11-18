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
	issues = ActiveTickets(issues, DefaultConfig())

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

func TestOpenTickets(t *testing.T) {
	config := DefaultConfig()
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Resolved = time.Now().AddDate(0, 0, -5)
	issue.Updated = time.Now()
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Resolved = time.Now().AddDate(0, -2, 0)
	issue.Updated = time.Now().AddDate(0, -1, -10)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "D"
	issue.Status = config.States.Closed
	issues = append(issues, issue)

	issues = OpenTickets(issues, config)

	if len(issues) != 1 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}
	if issues[0].Key != "B" {
		log.Println("TEST: wrong issue")
		t.Fail()
	}
}

func TestFilter(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issues = append(issues, issue)

	result := Filter(issues, func(issue *Issue) bool {
		return issue.Key == "A"
	})

	if len(result) != 1 {
		t.Fail()
	}
	if result[0].Key != "A" {
		t.Fail()
	}
}

func TestFixVersions(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.FixVersions = append(issue.FixVersions, "1", "2")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.FixVersions = append(issue.FixVersions, "3", "2")
	issues = append(issues, issue)

	versions := FixVersions(issues)

	str := ""
	for _, v := range versions {
		str += v
	}
	if str != "123" {
		t.Fail()
	}
}

func TestSecurityLevels(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.SecurityLevel = "A"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.SecurityLevel = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.SecurityLevel = "A"
	issues = append(issues, issue)

	security := SecurityLevels(issues)

	str := ""
	for _, v := range security {
		str += v
	}
	if str != "AB" {
		t.Fail()
	}
}

func TestTypes(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Type = "A"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Type = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Type = "A"
	issues = append(issues, issue)

	types := Types(issues)

	str := ""
	for _, v := range types {
		str += v
	}
	if str != "AB" {
		t.Fail()
	}
}

func TestLabels(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Labels = append(issue.Labels, "1", "2")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Labels = append(issue.Labels, "3", "2")
	issues = append(issues, issue)

	labels := Labels(issues)

	str := ""
	for _, v := range labels {
		str += v
	}
	if str != "123" {
		t.Fail()
	}
}

func TestComponents(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Components = append(issue.Components, "1", "2")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Components = append(issue.Components, "3", "2")
	issues = append(issues, issue)

	components := Components(issues)

	str := ""
	for _, v := range components {
		str += v
	}
	if str != "123" {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	list := []string{"A", "B", "C"}
	if !contains(list, "A") {
		t.Fail()
	}
	if contains(list, "D") {
		t.Fail()
	}
}

func TestFilterByFixVersion(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.FixVersions = append(issue.FixVersions, "A", "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.FixVersions = append(issue.FixVersions, "B")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.FixVersions = append(issue.FixVersions, "C", "D")
	issues = append(issues, issue)

	result := FilterByFixVersion(issues, "B")

	if len(result) != 2 {
		t.Fail()
	}
	keys := ""
	for _, i := range result {
		keys += i.Key
	}
	if keys != "AB" {
		t.Fail()
	}
}

func TestFilterBySecurityLevel(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.SecurityLevel = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.SecurityLevel = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.SecurityLevel = "D"
	issues = append(issues, issue)

	result := FilterBySecurityLevel(issues, "B")

	if len(result) != 2 {
		t.Fail()
	}
	keys := ""
	for _, i := range result {
		keys += i.Key
	}
	if keys != "AB" {
		t.Fail()
	}
}

func TestFilterByPriority(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Priority = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Priority = "B"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Priority = "D"
	issues = append(issues, issue)

	result := FilterByPriority(issues, "B")

	if len(result) != 2 {
		t.Fail()
	}
	keys := ""
	for _, i := range result {
		keys += i.Key
	}
	if keys != "AB" {
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

func TestFilterByLabel(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Components = append(issue.Components, "X")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Components = append(issue.Components, "Y", "X")
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Components = append(issue.Components, "Z")
	issues = append(issues, issue)

	// call function to test
	issues = FilterByComponent(issues, "X")

	if len(issues) != 2 {
		log.Println("TEST: issue count not expected")
		t.Fail()
	}

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}
	if keys != "AB" {

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

func TestOrderByDue(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Due = time.Now().AddDate(0, -1, 0)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Due = time.Now().AddDate(0, -1, -1)
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Due = time.Now().AddDate(0, 0, -20)
	issues = append(issues, issue)

	OrderByDue(issues)

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if keys != "BAC" {
		log.Println("TEST: wrong order")
		t.Fail()
	}
}

func TestOrderByPriority(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Priority = "b"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Priority = "a"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Priority = "c"
	issues = append(issues, issue)

	OrderByPriority(issues)

	keys := ""
	for _, i := range issues {
		keys += i.Key
	}

	if keys != "BAC" {
		log.Println("TEST: wrong order")
		t.Fail()
	}
}

func TestOrderByStatus(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.Status = "b"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.Status = "a"
	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"
	issue.Status = "c"
	issues = append(issues, issue)

	OrderByStatus(issues)

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
