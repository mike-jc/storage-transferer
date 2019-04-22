package system_test

import (
	"fmt"
	"service-recordingStorage/system"
	"testing"
	"time"
)

type addDurationFromStringData struct {
	OriginalTimeStr string
	OriginalTime    time.Time
	Duration        string
	ExpectedTimeStr string
	ExpectedTime    time.Time
}

func (d *addDurationFromStringData) InitTimes(location *time.Location) *addDurationFromStringData {
	var err error
	if d.OriginalTime, err = time.ParseInLocation("2006-01-02 15:04", d.OriginalTimeStr, location); err != nil {
		panic(fmt.Sprintf("Can not parse string %s, wrong time format: %s", d.OriginalTimeStr, err.Error()))
	}
	if d.ExpectedTime, err = time.ParseInLocation("2006-01-02 15:04", d.ExpectedTimeStr, location); err != nil {
		panic(fmt.Sprintf("Can not parse string %s, wrong time format: %s", d.ExpectedTimeStr, err.Error()))
	}
	return d
}

func TestAddDurationFromStringOK(t *testing.T) {

	testData := []addDurationFromStringData{
		{
			OriginalTimeStr: "2018-12-13 14:15",
			Duration:        "1 day",
			ExpectedTimeStr: "2018-12-14 14:15",
		},
		{
			OriginalTimeStr: "2018-12-13 14:15",
			Duration:        "5 days",
			ExpectedTimeStr: "2018-12-18 14:15",
		},
		{
			OriginalTimeStr: "2018-12-13 14:15",
			Duration:        "1 week 3 days",
			ExpectedTimeStr: "2018-12-23 14:15",
		},
		{
			OriginalTimeStr: "2018-12-13 14:15",
			Duration:        "1 year 10 months",
			ExpectedTimeStr: "2020-10-13 14:15",
		},
		{
			OriginalTimeStr: "2018-12-13 14:15",
			Duration:        "4 years 3 months 2 weeks 1 day",
			ExpectedTimeStr: "2023-03-28 14:15",
		},
	}

	for _, testRow := range testData {
		testRow.InitTimes(time.UTC)
		tResult := system.AddDurationFromString(testRow.OriginalTime, testRow.Duration)
		if !tResult.Equal(testRow.ExpectedTime) {
			t.Fatalf("Result time %+v is not equal to expected time %+v. Original time is %+v, duration is %s", tResult, testRow.ExpectedTime, testRow.OriginalTime, testRow.Duration)
		}
	}
}
