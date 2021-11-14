package ticketstats

import (
	"fmt"
	"log"
)

// SanitizeResult groups the sanitizer findings
type SanitizeResult struct {
	// issues with activity no set
	NoActivity []*Issue
	// invalid time bookings
	InvalidWorkLogs map[string][]WorkLog
}

func (issue *Issue) ExpectedActivity() string {
	activity := issue.CustomActivity
	if activity == "" {
		log.Println("WARNING: Activity not defined in ticket!", issue.Key)
		for _, l := range issue.LogWorks {
			if l.Activity != "" {
				activity = l.Activity
				break
			}
		}
	}

	if activity == "" {
		log.Println("WARNING: Na activity found for ticket!", issue.Key)
	}

	return activity
}

// AreBookingsValid checks if the work logs of the issue are consistent.
// The first value of the result is true if all logs are ok, the second
// is a list of invalid logs.
func (issue *Issue) AreBookingsValid() (bool, []WorkLog) {
	activity := issue.ExpectedActivity()
	valid := true
	invalidLogs := make([]WorkLog, 0)

	for _, l := range issue.LogWorks {
		if l.Activity == "" {
			log.Println("WARNING: Time booking without activity!", issue.Key,
				"expected activity", activity, "booking", l.ToString())
			valid = false
			invalidLogs = append(invalidLogs, l)
		}
		if l.Activity != activity {
			log.Println("WARNING: Time booking without wrong activity!",
				issue.Key, "expected activity", activity, "booking",
				l.ToString())
			valid = false
			invalidLogs = append(invalidLogs, l)
		}
	}

	return valid, invalidLogs
}

// Sanitize checks all issues for invalid state
func Sanitize(issues []*Issue) SanitizeResult {
	noActivity := make([]*Issue, 0)
	invalidLogs := make(map[string][]WorkLog)

	for _, issue := range issues {
		// Check if activity of ticket can be found
		if issue.ExpectedActivity() != "" {
			// Check tickets for wrong time bookings
			valid, logs := issue.AreBookingsValid()
			if !valid {
				invalidLogs[issue.Key] = logs
				fmt.Println("!!! Issue with wrong bookings:", issue.Key)
			}
		} else {
			noActivity = append(noActivity, issue)
		}
	}

	return SanitizeResult{
		NoActivity:      noActivity,
		InvalidWorkLogs: invalidLogs,
	}
}
