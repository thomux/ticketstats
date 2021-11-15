// CLI interface for ticketstats.
package main

import (
	"flag"

	"github.com/thomux/ticketstats/ticketstats"
)

func main() {
	var path string
	var project string
	var component string
	var jiraBase string
	var split bool

	flag.StringVar(&path, "csv", "JiraExport.csv", "path to Jira ticket export")
	flag.StringVar(&project, "project", "ICAS1", "Jira project key")
	flag.StringVar(&component, "component", "Linux", "Jira component name")
	flag.StringVar(&jiraBase, "jira", "", "Jira base URL")
	flag.BoolVar(&split, "splitByComponent", true, "split result by components")

	flag.Parse()

	ticketstats.Evaluate(path, project, component, jiraBase, true)
}
