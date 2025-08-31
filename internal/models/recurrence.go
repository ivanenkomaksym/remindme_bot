package models

import "time"

type Recurrence struct {
	Type        string         `json:"type"`         // e.g., "daily", "weekly", "monthly", "interval", "custom"
	Interval    int            `json:"interval"`     // e.g., every N days/hours/minutes
	Weekdays    []time.Weekday `json:"weekdays"`     // For weekly recurrence (e.g., [Tuesday, Thursday])
	DayOfMonth  int            `json:"day_of_month"` // For monthly recurrence
	TimeOfDay   string         `json:"time_of_day"`  // e.g., "14:00"
	StartDate   *time.Time     `json:"start_date"`   // When recurrence begins
	EndDate     *time.Time     `json:"end_date"`     // When recurrence ends (optional)
	Occurrences *int           `json:"occurrences"`  // Number of times to repeat (optional)
	CustomDays  []int          `json:"custom_days"`  // For custom patterns (e.g., every 5 days)
}

func (r *Recurrence) IsWeekly() bool {
	return r.Type == "weekly" && len(r.Weekdays) > 0
}

func (r *Recurrence) IsInterval() bool {
	return r.Type == "interval" && r.Interval > 0
}

func (r *Recurrence) IsMonthly() bool {
	return r.Type == "monthly" && r.DayOfMonth > 0
}

func (r *Recurrence) IsDaily() bool {
	return r.Type == "daily"
}
