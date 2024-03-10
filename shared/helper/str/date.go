package str

import (
	"strings"
	"time"
)

func TrimWhiteSpace(inputString string) string {
	return strings.TrimSpace(inputString)
}

func ChangeStringToDateTime(dateInString string) (*time.Time, string) {

	if dateInString == "" {

		errorResult := "Please input your date"
		return nil, errorResult
	}

	resultDate, err := time.Parse("2006-01-02", TrimWhiteSpace(dateInString))
	if err != nil {

		errorResult := "Please input your date with format \"yyyy-MM-dd\"" + err.Error()
		return nil, errorResult
	}

	return &resultDate, ""
}

func GetCurrentDateAndZeroTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return currentDate
}

func GetCurrentDateTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	return now
}
