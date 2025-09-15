package entities

import "errors"

type RecurrenceType int64

const (
	Once RecurrenceType = iota
	Daily
	Weekly
	Monthly
	Interval
	Custom
)

var RecurrenceTypeValues = []RecurrenceType{
	Once,
	Daily,
	Weekly,
	Monthly,
	Interval,
	Custom,
}

func (r RecurrenceType) String() string {
	switch r {
	case Once:
		return "Once"
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
	case "Once":
		return Once, nil
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
		return 0, errors.New("invalid recurrence type")
	}
}
