package ticketstats

import (
	"log"
	"testing"
	"time"
)

func TestAreBookingsValid(t *testing.T) {
	config := DefaultConfig()
	issue := NewIssue()
	issue.Key = "A"
	issue.CustomActivity = "123457"

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123457",
		Hours:    2.5,
		Date:     time.Now(),
	})

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123457",
		Hours:    3.5,
		Date:     time.Now(),
	})

	valid, _ := issue.AreBookingsValid(false, config)

	if !valid {
		log.Println("TEST: issue should be valid")
		t.Fail()
	}

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123456",
		Hours:    1.5,
		Date:     time.Now(),
	})

	valid, logs := issue.AreBookingsValid(false, config)

	if valid {
		log.Println("TEST: issue should be invalid")
		t.Fail()
	}
	if len(logs) != 1 {
		log.Println("TEST: wrong log count")
		t.Fail()
	}
	if logs[0].Activity != "123456" {
		log.Println("TEST: wrong log")
		t.Fail()
	}
}

func TestSanitize(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.Key = "A"
	issue.CustomActivity = "123457"

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123457",
		Hours:    2.5,
		Date:     time.Now(),
	})

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123457",
		Hours:    3.5,
		Date:     time.Now(),
	})

	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "B"
	issue.CustomActivity = "123457"

	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Activity: "123456",
		Hours:    1.5,
		Date:     time.Now(),
	})

	issues = append(issues, issue)

	issue = NewIssue()
	issue.Key = "C"

	issues = append(issues, issue)

	result := Sanitize(issues, false, DefaultConfig())

	if len(result.NoActivity) != 1 {
		log.Println("TEST: wrong count of issues with no activity")
		t.Fail()
	}
	if result.NoActivity[0].Key != "C" {
		log.Println("TEST: wrong issue with no activity")
		t.Fail()
	}

	if len(result.InvalidWorkLogs) != 1 {
		log.Println("TEST: wrong count of issues with invalid bookings")
		t.Fail()
	}
	il := result.InvalidWorkLogs[0]
	if il.Issue.Key != "B" {
		log.Println("TEST: wrong issue with invalid bookings")
		t.Fail()
	}
	if il.Logs[0].Activity != "123456" {
		log.Println("TEST: wrong invalid work log")
		t.Fail()
	}
}

func TestIgnoreOld(t *testing.T) {
	t.Fail()
}
