package ticketstats

import (
	"log"
	"testing"
)

func TestConvertWork(t *testing.T) {
	if convertWork("5400") != 1.5 {
		t.Fail()
	}
}

func TestConvertWorkLog(t *testing.T) {
	// Log Work = [entry description (containing ExAc as ExecutionActivity:XXXXXX)];[Date];[user (always asw_qm_service)];[time spend as seconds]
	data := "Hallo Welt\nExecutionActivity:123457\nmore text;31/Aug/21 12:57 PM;auser;45000"
	work := convertWorkLog(data)

	if work.Activity != "123457" {
		log.Println("TEST: activity wrong", work.ToString())
		t.Fail()
	}

	if work.Hours != 12.5 {
		log.Println("TEST: hours", work.ToString())
		t.Fail()
	}

	if work.Date != convertDate("31/Aug/21 12:57 PM") {
		log.Println("TEST: date", work.ToString())
		t.Fail()
	}
}

func TestConvertDate(t *testing.T) {
	date := convertDate("31/Aug/21 3:57 PM")
	if date.Year() != 2021 {
		log.Println("TEST: year", date)
		t.Fail()
	}

	if date.Month() != 8 {
		log.Println("TEST: month", date)
		t.Fail()
	}

	if date.Day() != 31 {
		log.Println("TEST: day", date)
		t.Fail()
	}

	if date.Hour() != 15 {
		log.Println("TEST: hour", date)
		t.Fail()
	}

	if date.Minute() != 57 {
		log.Println("TEST: minute", date)
		t.Fail()
	}
}
