package modelsDataInstance

import "time"

type Meeting struct {
	Id          int64       `json:"id"`
	Type        MeetingType `json:"type"`
	User        User        `json:"user"`
	Guest       Guest       `json:"guest"`
	DateUTC     time.Time   `json:"date"`
	Timezone    string      `json:"timezone"`
	Description string      `json:"description"`
}

func (m *Meeting) DateInTz() time.Time {
	if location, err := time.LoadLocation(m.Timezone); err == nil {
		return m.DateUTC.In(location)
	} else {
		return m.DateUTC
	}
}
