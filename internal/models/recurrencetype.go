package models

type RecurrenceType int64

const (
	Daily RecurrenceType = iota
	Weekly
	Monthly
	Interval
	Custom
)

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
