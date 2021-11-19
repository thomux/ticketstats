package ticketstats

import (
	"log"
	"strings"
	"time"
)

// SanitizeResult groups the sanitizer findings
type SanitizeResult struct {
	// issues with activity no set
	NoActivity []*Issue
	// invalid time bookings
	InvalidWorkLogs []InvalidWorkLog
}

// InvalidWorkLog groups all data for an invalid work log.
type InvalidWorkLog struct {
	Issue *Issue
	Logs  []WorkLog
}

// NewInvalidWorkLog initializes a new InvalidWorkLog.
func NewInvalidWorkLog(issue *Issue) InvalidWorkLog {
	var invalidLog InvalidWorkLog
	invalidLog.Issue = issue
	invalidLog.Logs = make([]WorkLog, 0)
	return invalidLog
}

// AreBookingsValid checks if the work logs of the issue are consistent.
// The first value of the result is true if all logs are ok, the second
// is a list of invalid logs.
func (issue *Issue) AreBookingsValid(ignoreOld bool,
	config Config) (bool, []WorkLog) {

	activity := strings.TrimSpace(issue.CustomActivity)
	valid := true
	invalidLogs := make([]WorkLog, 0)

	if activity == "" {
		return true, invalidLogs
	}

	start := time.Now().AddDate(0, 0, -time.Now().Day()-1)
	for _, l := range issue.LogWorks {
		if ignoreOld && l.Date.Before(start) {
			continue
		}
		if l.Activity == "" {
			log.Println("ERROR: WorkLog without activity!",
				l.ToString(config), issue.Key)
		}

		if strings.Compare(l.Activity, activity) != 0 {
			valid = false
			invalidLogs = append(invalidLogs, l)
		}
	}

	return valid, invalidLogs
}

// Sanitize checks all issues for invalid state
func Sanitize(issues []*Issue, ignoreOld bool, config Config) SanitizeResult {
	noActivity := make([]*Issue, 0)
	invalidLogs := make([]InvalidWorkLog, 0)

	for _, issue := range issues {
		// Check if activity of ticket can be found
		if issue.CustomActivity != "" {
			// Check tickets for wrong time bookings
			valid, logs := issue.AreBookingsValid(ignoreOld, config)
			if !valid {
				invalidLog := NewInvalidWorkLog(issue)
				invalidLog.Logs = append(invalidLog.Logs, logs...)
				invalidLogs = append(invalidLogs, invalidLog)
			}
		} else {
			if issue.Status != config.States.Closed {
				noActivity = append(noActivity, issue)
			}
		}
	}

	return SanitizeResult{
		NoActivity:      noActivity,
		InvalidWorkLogs: invalidLogs,
	}
}
