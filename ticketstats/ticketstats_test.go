package ticketstats

import (
	"testing"
	"time"
)

func TestCalcHours(t *testing.T) {
	issues := make([]*Issue, 0)

	issue := NewIssue()
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    10.0,
		Date:     time.Now().AddDate(0, 0, -2),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    20.0,
		Date:     time.Now().AddDate(0, 0, -10),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    30.0,
		Date:     time.Now().AddDate(0, -1, 0),
		Activity: "123456",
	})
	issue.LogWorks = append(issue.LogWorks, WorkLog{
		Hours:    40.0,
		Date:     time.Now().AddDate(0, -6, 0),
		Activity: "123456",
	})
	issues = append(issues, issue)

	work := calcHours(issues)

	if work[0] != 10.0 ||
		work[1] != 30.0 ||
		work[2] != 60.0 ||
		work[3] != 100.0 {
		t.Fail()
	}
}

func TestCalcFTE(t *testing.T) {
	work := []Work{80.0, 170.0, 510.0, 1040.0}
	fte := calcFTE(work)
	if fte[0] != 2.0 ||
		fte[1] != 1.0 ||
		fte[2] != 1.0 ||
		fte[3] != 0.5 {
		t.Fail()
	}
}
