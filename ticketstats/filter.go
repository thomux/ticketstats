package ticketstats

import (
	"sort"
	"strings"
	"time"
)

// ActiveTickets returns all active tickets.
func ActiveTickets(issues []*Issue) []*Issue {
	active := make([]*Issue, 0)
	noDate := time.Time{}

	for _, issue := range issues {
		if issue.Resolved == noDate ||
			lastMonth(issue.Updated) {
			active = append(active, issue)
		}
	}

	return active
}

// FilterByProject only returns the tickets matching the given project key.
func FilterByProject(issues []*Issue, project string) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if strings.HasPrefix(issue.Key, project) {
			res = append(res, issue)
		}
	}

	return res
}

// CreatedLastWeek returns all issues created during the last 7 days.
func CreatedLastWeek(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastWeek(issue.Created) {
			res = append(res, issue)
		}
	}

	return res
}

// CreatedLastMonth returns all issues created during the last month.
func CreatedLastMonth(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastMonth(issue.Created) {
			res = append(res, issue)
		}
	}

	return res
}

// CreatedLastQuarter returns all issues created during the last three month.
func CreatedLastQuarter(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastQuarter(issue.Created) {
			res = append(res, issue)
		}
	}

	return res
}

// CreatedLastYear returns all issues created during the last year
func CreatedLastYear(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastYear(issue.Created) {
			res = append(res, issue)
		}
	}

	return res
}

// ClosedLastWeek returns all issues resolved during the last week.
func ClosedLastWeek(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastWeek(issue.Resolved) {
			res = append(res, issue)
		}
	}

	return res
}

// ClosedLastMonth returns all issues resolved during the last month.
func ClosedLastMonth(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastMonth(issue.Resolved) {
			res = append(res, issue)
		}
	}

	return res
}

// ClosedLastQuarter returns all issues resolved during the last three month.
func ClosedLastQuarter(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastQuarter(issue.Resolved) {
			res = append(res, issue)
		}
	}

	return res
}

// ClosedLastYear returns all issues resolved during the last year.
func ClosedLastYear(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if lastYear(issue.Resolved) {
			res = append(res, issue)
		}
	}

	return res
}

// OlderThanOneMonth returns all tickets older than one month.
func OlderThanOneMonth(issues []*Issue) []*Issue {
	res := make([]*Issue, 0)
	lm := time.Now().AddDate(0, -1, 0)

	for _, issue := range issues {
		if issue.Created.Before(lm) {
			res = append(res, issue)
		}
	}

	return res
}

// FilterByType returns all issues matching the given type.
func FilterByType(issues []*Issue, ticketType string) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		if issue.Type == ticketType {
			res = append(res, issue)
		}
	}

	return res
}

// FilterByComponent returns all issues matching the given component.
func FilterByComponent(issues []*Issue, component string) []*Issue {
	res := make([]*Issue, 0)

	for _, issue := range issues {
		for _, c := range issue.Components {
			if c == component {
				res = append(res, issue)
			}
		}
	}

	return res
}

// OrderByCreated orders the issues by created date.
func OrderByCreated(issues []*Issue) {
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Created.Before(issues[j].Created)
	})
}

// Age returns the difference from now to the given date as days.
func Age(date time.Time) int {
	diff := time.Since(date)
	return int(diff.Hours() / 24)
}

// IssuesByComponent groups the given issues by component.
func IssuesByComponent(issues []*Issue) map[string][]*Issue {
	components := make(map[string][]*Issue)

	for _, issue := range issues {
		for _, component := range issue.Components {
			list, ok := components[component]
			if !ok {
				list = make([]*Issue, 0)
			}
			list = append(list, issue)
			components[component] = list
		}
	}

	return components
}

// lastWeek creates a date 7 days back from now.
func lastWeek(date time.Time) bool {
	lm := time.Now().AddDate(0, 0, -7)
	return date.After(lm)
}

// lastMonth creates a date one month back from now.
func lastMonth(date time.Time) bool {
	lm := time.Now().AddDate(0, -1, 0)
	return date.After(lm)
}

// lastQuarter creates a date three month back from now.
func lastQuarter(date time.Time) bool {
	lm := time.Now().AddDate(0, -3, 0)
	return date.After(lm)
}

// lastYear creates a date one year back from now.
func lastYear(date time.Time) bool {
	lm := time.Now().AddDate(-1, 0, 0)
	return date.After(lm)
}
