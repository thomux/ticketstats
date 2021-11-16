package ticketstats

import (
	"sort"
	"strings"
	"time"
)

// ActiveTickets returns all active tickets.
func ActiveTickets(issues []*Issue) []*Issue {
	noDate := time.Time{}

	return Filter(issues, func(issue *Issue) bool {
		return (issue.Resolved == noDate && issue.Status != "Closed") || lastMonth(issue.Updated)
	})
}

func OpenTickets(issues []*Issue) []*Issue {
	noDate := time.Time{}

	return Filter(issues, func(issue *Issue) bool {
		return issue.Resolved == noDate && issue.Status != "Closed"
	})
}

func Filter(issues []*Issue, test func(issue *Issue) bool) []*Issue {
	result := make([]*Issue, 0)

	for _, issue := range issues {
		if test(issue) {
			result = append(result, issue)
		}
	}

	return result
}

func FixVersions(issues []*Issue) []string {
	result := make([]string, 0)
	fixVersions := make(map[string]int)

	for _, issue := range issues {
		for _, version := range issue.FixVersions {
			_, ok := fixVersions[version]
			if !ok {
				fixVersions[version] = 1
				result = append(result, version)
			}
		}
	}

	return result
}

func SecurityLevels(issues []*Issue) []string {
	result := make([]string, 0)
	security := make(map[string]int)

	for _, issue := range issues {
		_, ok := security[issue.SecurityLevel]
		if !ok {
			security[issue.SecurityLevel] = 1
			result = append(result, issue.SecurityLevel)
		}
	}

	return result
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func FilterByFixVersion(issues []*Issue, fixVersion string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return contains(issue.FixVersions, fixVersion)
	})
}

func FilterBySecurityLevel(issues []*Issue, securityLevel string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return issue.SecurityLevel == securityLevel
	})
}

func FilterByPriority(issues []*Issue, priority string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return issue.Priority == priority
	})
}

// FilterByProject only returns the tickets matching the given project key.
func FilterByProject(issues []*Issue, project string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return strings.HasPrefix(issue.Key, project)
	})
}

// CreatedLastWeek returns all issues created during the last 7 days.
func CreatedLastWeek(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastWeek(issue.Created)
	})
}

// CreatedLastMonth returns all issues created during the last month.
func CreatedLastMonth(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastMonth(issue.Created)
	})
}

// CreatedLastQuarter returns all issues created during the last three month.
func CreatedLastQuarter(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastQuarter(issue.Created)
	})
}

// CreatedLastYear returns all issues created during the last year
func CreatedLastYear(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastYear(issue.Created)
	})
}

// ClosedLastWeek returns all issues resolved during the last week.
func ClosedLastWeek(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastWeek(issue.Resolved)
	})
}

// ClosedLastMonth returns all issues resolved during the last month.
func ClosedLastMonth(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastMonth(issue.Resolved)
	})
}

// ClosedLastQuarter returns all issues resolved during the last three month.
func ClosedLastQuarter(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastQuarter(issue.Resolved)
	})
}

// ClosedLastYear returns all issues resolved during the last year.
func ClosedLastYear(issues []*Issue) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return lastYear(issue.Resolved)
	})
}

// OlderThanOneMonth returns all tickets older than one month.
func OlderThanOneMonth(issues []*Issue) []*Issue {
	lm := time.Now().AddDate(0, -1, 0)
	return Filter(issues, func(issue *Issue) bool {
		return issue.Created.Before(lm)
	})
}

// FilterByType returns all issues matching the given type.
func FilterByType(issues []*Issue, ticketType string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return issue.Type == ticketType
	})
}

// FilterByComponent returns all issues matching the given component.
func FilterByComponent(issues []*Issue, component string) []*Issue {
	return Filter(issues, func(issue *Issue) bool {
		return contains(issue.Components, component)
	})
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
