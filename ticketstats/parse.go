package ticketstats

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func convertWork(data string) Work {
	secs, err := strconv.Atoi(data)
	if err != nil {
		log.Println("ERROR:", err)
		return 0
	} else {
		return Work(float64(secs) / 3600)
	}
}

func convertWorkLog(data string) WorkLog {
	// Log Work = [entry description (containing ExAc as ExecutionActivity:XXXXXX)];[Date];[user (always asw_qm_service)];[time spend as seconds]
	var hours Work
	var date time.Time
	var exAc string

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ExecutionActivity:") {
			exAc = line[18:]
		}
	}

	line := lines[len(lines)-1]
	i := strings.LastIndex(line, ";")
	hours = convertWork(line[i+1:])

	// skip user
	tmp := line[:i]
	i = strings.LastIndex(tmp, ";")
	// get date
	tmp = tmp[:i]
	i = strings.LastIndex(tmp, ";")
	date = convertDate(tmp[i+1:])

	return WorkLog{
		Hours:    hours,
		Date:     date,
		Activity: exAc,
	}
}

func convertDate(data string) time.Time {
	layout := "02/Jan/06 3:04 PM"
	t, err := time.Parse(layout, data)
	if err != nil {
		log.Println("ERROR:", err)
		return time.Time{}
	}
	return t
}

func Parse(path string) []*Issue {
	records := readCsvFile("JiraExport.csv")

	header := records[0]
	data := records[1:]
	issues := make([]*Issue, 0)

	for _, d := range data {
		issue := NewIssue()
		for i, val := range d {
			val = strings.TrimSpace(val)

			// skip empty values
			if val == "" {
				continue
			}

			key := header[i]

			switch key {
			case "Summary":
				issue.Summary = val
			case "Issue key":
				issue.Key = val
			case "Issue id":
				issue.Id = val
			case "Parent id":
				issue.Parent = val
			case "Issue Type":
				issue.Type = val
			case "Status":
				issue.Status = val
			case "Priority":
				issue.Priority = val
			case "Assignee":
				issue.Assignee = val
			case "Creator":
				issue.Creator = val
			case "Created":
				issue.Created = convertDate(val)
			case "Updated":
				issue.Updated = convertDate(val)
			case "Last Viewed":
				issue.LastViewed = convertDate(val)
			case "Affects Version/s":
				issue.AffectsVersions = append(issue.AffectsVersions, val)
			case "Fix Version/s":
				issue.FixVersions = append(issue.FixVersions, val)
			case "Component/s":
				issue.Components = append(issue.Components, val)
			case "Log Work":
				issue.LogWorks = append(issue.LogWorks, convertWorkLog(val))
			case "Original Estimate":
				issue.OriginalEstimate = convertWork(val)
			case "Remaining Estimate":
				issue.RemainingEstimate = convertWork(val)
			case "Time Spent":
				issue.TimeSpend = convertWork(val)
			case "Σ Original Estimate":
				issue.SumOriginalEstimate = convertWork(val)
			case "Σ Remaining Estimate":
				issue.SumRemainingEstimate = convertWork(val)
			case "Σ Time Spent":
				issue.SumTimeSpend = convertWork(val)
			case "Security Level":
				issue.SecurityLevel = val
			case "Labels":
				issue.Labels = append(issue.Labels, val)
			case "Resolution":
				issue.Resolution = val
			case "Resolved":
				issue.Resolved = convertDate(val)
			case "Due Date":
				issue.Due = convertDate(val)
			case "Outward issue link (Blocks)":
				issue.LinkBlocks = append(issue.LinkBlocks, val)
			case "Outward issue link (Causes)":
				issue.LinkCauses = append(issue.LinkCauses, val)
			case "Outward issue link (Cloners)":
				issue.LinkCloners = append(issue.LinkCloners, val)
			case "Outward issue link (Dependency)":
				issue.LinkDependencies = append(issue.LinkDependencies, val)
			case "Outward issue link (Duplicate)":
				issue.LinkDuplicates = append(issue.LinkDuplicates, val)
			case "Outward issue link (Issue split)":
				issue.LinkIssueSplits = append(issue.LinkIssueSplits, val)
			case "Outward issue link (Part)":
				issue.LinkParts = append(issue.LinkParts, val)
			case "Outward issue link (Relates)":
				issue.LinkRelates = append(issue.LinkRelates, val)
			case "Outward issue link (Relation)":
				issue.LinkRelations = append(issue.LinkRelations, val)
			case "Outward issue link (Triggers)":
				issue.LinkTriggers = append(issue.LinkTriggers, val)
			case "Outward issue link (linkIssue)":
				issue.LinkLinkIssues = append(issue.LinkLinkIssues, val)
			case "Outward issue link (parent)":
				issue.LinkParents = append(issue.LinkParents, val)
			case "Custom field (External ID)":
				issue.CustomExternalId = val
			case "Custom field (Supplier reference)":
				issue.CustomSupplierRef = val
			case "Custom field (ICAS Variant)":
				issue.CustomVariant = val
			case "Custom field (Booking Account)":
				issue.CustomActivity = val
			case "Custom field (Bug-Category)":
				issue.CustomCategory = val
			}
		}
		issues = append(issues, issue)
	}
	return issues
}
