package modelsDataStorage

import "fmt"

type State struct {
	Range LoadRange
}

type LoadRange struct {
	Start int64
	End   int64
	Limit int64
}

func (r *LoadRange) RangeHeader() string {
	if r.End > 0 {
		return fmt.Sprintf("bytes=%d-%d", r.Start, r.End)
	} else {
		return ""
	}
}

func (r *LoadRange) ContentRangeHeader() string {
	if r.End > 0 && r.Limit > 0 {
		return fmt.Sprintf("bytes %d-%d/%d", r.Start, r.End, r.Limit)
	} else {
		return ""
	}
}
