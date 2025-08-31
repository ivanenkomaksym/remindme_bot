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
		return "daily"
	case Weekly:
		return "weekly"
	case Monthly:
		return "monthly"
	case Interval:
		return "interval"
	case Custom:
		return "custom"
	default:
		return "unknown"
	}
}

func ToRecurrenceType(s string) (RecurrenceType, error) {
	switch s {
	case "daily":
		return Daily, nil
	case "weekly":
		return Weekly, nil
	case "monthly":
		return Monthly, nil
	case "interval":
		return Interval, nil
	case "custom":
		return Custom, nil
	default:
		return -1, errors.New("unknown recurrence")
	}
}
