package ticketstats

import "log"

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

func (issue *Issue) AreBookingsValid() bool {
	activity := issue.ExpectedActivity()
	valid := true

	for _, l := range issue.LogWorks {
		if l.Activity == "" {
			log.Println("WARNING: Time booking without activity!", issue.Key,
				"expected activity", activity, "booking", l.ToString())
			valid = false
		}
		if l.Activity != activity {
			log.Println("WARNING: Time booking without wrong activity!", issue.Key,
				"expected activity", activity, "booking", l.ToString())
			valid = false
		}
	}

	return valid
}
