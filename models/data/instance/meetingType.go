package modelsDataInstance

import "fmt"

type MeetingType struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

func (t *MeetingType) NameAndDuration() string {
	return fmt.Sprintf("%s (%d min)", t.Name, t.Duration)
}
