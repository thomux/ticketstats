// ticketstats generates statistics and further report for a set of csv exported
// Jira tickets.
package ticketstats

import (
	"fmt"
	"log"
)

// Evaluate generates a full report for the exported tickets
func Evaluate(path string,
	project string,
	component string,
	splitByComponent bool) {

	// read issues form csv
	issues := Parse(path)

	if project != "" {
		issues = FilterByProject(issues, project)
	}
	if component != "" {
		issues = FilterByComponent(issues, component)
	}

	// Reduce to active tickets
	active := ActiveTickets(issues)
	log.Println("INFO:", len(active), "active tickets.")

	// Check tickets for issues
	Sanitize(active)

	oldBugs := OldBugs(active)
	log.Println("INFO:", len(oldBugs), "old bug tickets.")
	OrderByCreated(oldBugs)
	str := "INFO:top 10 oldest old bug tickets:\n"
	for _, issue := range oldBugs[:10] {
		str += fmt.Sprintf("- %s %s (age: %d days)\n",
			issue.Key,
			issue.Summary,
			Age(issue.Created))
	}
	log.Println(str)

	for tp, tr := range ResultionTimesByType(issues) {
		log.Printf("INFO: statistics for ticket type %s:\n%s", tp, tr.ToString())
	}
}
