package entities

import (
	"bytes"
	"encoding/json"
	"errors"
)

type RecurrenceType int64

const (
	Once RecurrenceType = iota
	Daily
	Weekly
	Monthly
	Interval
)

var RecurrenceTypeValues = []RecurrenceType{
	Once,
	Daily,
	Weekly,
	Monthly,
	Interval,
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
	default:
		return 0, errors.New("invalid recurrence type")
	}
}

// MarshalJSON marshals the enum as a quoted json string
func (r RecurrenceType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *RecurrenceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*r, _ = ToRecurrenceType(s)
	return nil
}
