package system

import (
	"regexp"
	"strconv"
	"time"
)

func AddDurationFromString(t time.Time, s string) time.Time {

	parseTime := func(s string, part string) int {
		amount := 0
		r := regexp.MustCompile(`(\d+)(\s+)` + part + `(s?)`)

		if matches := r.FindAllStringSubmatch(s, -1); matches != nil {
			for _, subMatches := range matches {
				if v, err := strconv.Atoi(subMatches[1]); err == nil {
					amount += v
				}
			}
		}
		return amount
	}

	years := parseTime(s, "year")
	months := parseTime(s, "month")
	weeks := parseTime(s, "week")
	days := parseTime(s, "day")

	return t.AddDate(years, months, days+7*weeks)
}
