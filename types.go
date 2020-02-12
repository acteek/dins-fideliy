package main

import "time"

//Action for Subscripton
type Action int

// Subscription can be Create or Delete
const (
	Create Action = iota
	Delete
)

// Subscription data fields
type Subscription struct {
	ChatID int64
	Action Action
}

// TimeRange time range for publish subs
type TimeRange struct {
	Start string
	End   string
}

func parseTime(str string) time.Time {
	time, _ := time.Parse("15:04", str)
	return time
}

func (r *TimeRange) contain(time string) bool {
	start := parseTime(r.Start)
	end := parseTime(r.End)
	check := parseTime(time)

	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

