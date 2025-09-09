package models

import "errors"

type RecurrenceType int64

const (
	Daily RecurrenceType = iota
	Weekly
	Monthly
	Interval
	Custom
)

var RecurrenceTypeValues = []RecurrenceType{
	Daily,
	Weekly,
	Monthly,
	Interval,
	Custom,
}

func (r RecurrenceType) String() string {
	switch r {
	case Daily:
		return "Daily"
	case Weekly:
		return "Weekly"
	case Monthly:
		return "Monthly"
	case Interval:
		return "Interval"
	case Custom:
		return "Custom"
	default:
		return "unknown"
	}
}

func ToRecurrenceType(s string) (RecurrenceType, error) {
	switch s {
	case "Daily":
		return Daily, nil
	case "Weekly":
		return Weekly, nil
	case "Monthly":
		return Monthly, nil
	case "Interval":
		return Interval, nil
	case "Custom":
		return Custom, nil
	default:
		return -1, errors.New("unknown recurrence")
	}
}
