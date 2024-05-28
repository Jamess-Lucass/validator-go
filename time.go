package schema

import "time"

type TimeSchema struct {
	Schema[time.Time]
}

var _ ISchema = (*TimeSchema)(nil)

func Time() *TimeSchema {
	return &TimeSchema{}
}
