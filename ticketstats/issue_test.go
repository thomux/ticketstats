package ticketstats

import (
	"log"
	"testing"
	"time"
)

func TestFormatWork(t *testing.T) {
	var work Work = 7.5
	if formatWork(work) != "7.50h" {
		log.Println("TEST: 7.5h wrong:", formatWork(work))
		t.Fail()
	}

	work = 12.5
	if formatWork(work) != "1d 4.50h" {
		log.Println("TEST: 12.5 wrong:", formatWork(work))
		t.Fail()
	}

	work = 60.5
	if formatWork(work) != "1w 2d 4.50h" {
		log.Println("TEST: 60.5 wrong:", formatWork(work))
		t.Fail()
	}
}

func TestIsResolved(t *testing.T) {
	issue := NewIssue()
	issue.Resolved = time.Now()

	if !issue.IsResolved() {
		log.Println("TEST: resolved issue wrong")
		t.Fail()
	}

	issue = NewIssue()

	if issue.IsResolved() {
		log.Println("TEST: not resolved issue wrong")
		t.Fail()
	}
}
