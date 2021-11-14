package main

import (
	"fmt"

	"github.com/thomux/ticketstats/ticketstats"
)

func main() {
	issues := ticketstats.Parse("JiraExport.csv")

	for _, issue := range issues[:3] {
		fmt.Println(issue.ToString())
	}

	for _, issue := range issues {
		if issue.ExpectedActivity() == "" {
			continue
		}
		if !issue.AreBookingsValid() {
			fmt.Println("!!! Issue with wrong bookings:", issue.Key)
		}
	}
}
