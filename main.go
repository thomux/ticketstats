package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Issue fields

Summary = string
Issue key = string
Issue id = string
Parent id = string
Issue Type = string
Status = string
Priority = string
Assignee = string
Creator = string
Created = date 20/Oct/21 2:24 PM
Updated = date 13/Nov/21 9:23 AM
Last Viewed = date 13/Nov/21 9:23 AM
Affects Version/s = string list 4.5.1
Fix Version/s = string list 4.7.1
Component/s = string list Linux
Log Work = workLog list 8.333333333333334 on ExAc 476355
Log Work = [entry description (containing ExAc as ExecutionActivity:XXXXXX)];[Date];[user (always asw_qm_service)];[time spend as seconds]
Original Estimate = float64, convert form seconds to hours 0
Remaining Estimate = float64, convert form seconds to hours
Time Spent = float64, convert form seconds to hours 318240
Σ Original Estimate = float64, convert form seconds to hours 0
Σ Remaining Estimate = float64, convert form seconds to hours
Σ Time Spent = float64, convert form seconds to hours 318240
Security Level = string
Affects Version/s = string
Component/s = string list Linux
Labels = string list IMPL_APPROVED
Resolution = string
Resolved = date
Due Date = date

Outward issue link (Blocks) = string  # ignore for clustering
Outward issue link (Causes) = string  # ignore for clustering
Outward issue link (Cloners) = string # cluster, order by creation date
Outward issue link (Dependency) = string # ignore for clustering
Outward issue link (Duplicate) = string  # cluster, order by creation date
Outward issue link (Issue split) = string # cluster, order by creation date
Outward issue link (Part) = string # cluster, order by creation date
Outward issue link (Relates) = string # ignore for clustering
Outward issue link (Relation) = string # ignore for clustering
Outward issue link (Triggers) = string # ignore for clustering
Outward issue link (linkIssue) = string # ignore for clustering
Outward issue link (parent) = string # cluster, order by link

Custom field (External ID) = string
Custom field (Supplier reference) = string
Custom field (ICAS Variant) = string
Custom field (Booking Account) = string
Custom field (Bug-Category) = string

Childs = Issue list

*/

type Work float64

type WorkLog struct {
	Hours    Work
	Activity string
}

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

func convertTime(data string) (Work, error) {
	secs, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	} else {
		return Work(float64(secs) / 3600), nil
	}
}

func convertWorkLog(data string) WorkLog {
	var err error
	var hours Work
	exAc := ""
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ExecutionActivity:") {
			exAc = line[18:]
			if err != nil {
				log.Println("ERROR:", err)
			}
		}

		if strings.HasPrefix(line, "Info=") {
			i := strings.LastIndex(line, ";")
			hours, err = convertTime(line[i+1:])
			if err != nil {
				log.Println("ERROR:", err)
			}
		}
	}

	return WorkLog{
		Hours:    hours,
		Activity: exAc,
	}
}

func convertDate(data string) (time.Time, error) {
	layout := "02/Jan/06 3:04 PM"
	return time.Parse(layout, data)
}

func main() {
	records := readCsvFile("JiraExport.csv")

	header := records[0]
	data := records[1:]

	for _, d := range data[:10] {

		for i, v := range d {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}

			if header[i] == "Attachment" || header[i] == "Comment" || header[i] == "Description" || header[i] == "Custom field (External Description)" {
				continue
			}

			if header[i] == "Log Work" {
				workLog := convertWorkLog(v)
				fmt.Printf("%s = %+v\n", header[i], workLog)
				continue
			}

			date, err := convertDate(v)
			if err == nil {
				fmt.Println(header[i], "=", date)
				continue
			}

			if strings.Contains(header[i], "Estimate") || strings.Contains(header[i], "Time") {
				work, err := convertTime(v)
				if err != nil {
					log.Println("ERROR:", err)
					continue
				}
				fmt.Printf("%s = %fh\n", header[i], work)
				continue
			}

			fmt.Println(header[i], "=", v)
		}
		fmt.Println("--------------------")
	}
}
