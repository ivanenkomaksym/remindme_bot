package models

import "time"

type Recurrence struct {
	Type       string         `json:"type"`          // e.g., "daily", "weekly", "monthly", "interval", "custom"
	Interval   int            `json:"interval"`      // e.g., every N days/hours/minutes
	Weekdays   []time.Weekday `json:"weekdays"`      // For weekly recurrence (e.g., [Tuesday, Thursday])
	DayOfMonth []int          `json:"days_of_month"` // For monthly recurrence (e.g., [1, 15])
	TimeOfDay  string         `json:"time_of_day"`   // e.g., "14:00"
	StartDate  *time.Time     `json:"start_date"`    // When recurrence begins
	EndDate    *time.Time     `json:"end_date"`      // When recurrence ends (optional)
}

// WithStartDate adds a start date to a recurrence
func (r *Recurrence) WithStartDate(startDate time.Time) *Recurrence {
	r.StartDate = &startDate
	return r
}

// WithEndDate adds an end date to a recurrence
func (r *Recurrence) WithEndDate(endDate time.Time) *Recurrence {
	r.EndDate = &endDate
	return r
}

// CustomWeekly creates a custom weekly recurrence on specific weekdays
func CustomWeekly(weekdays []time.Weekday, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "weekly",
		Weekdays:  weekdays,
		TimeOfDay: timeOfDay,
	}
}

// DailyAt creates a daily recurrence at a specific time
func DailyAt(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "daily",
		TimeOfDay: timeOfDay,
	}
}

// MonthlyOnDay creates a monthly recurrence on a specific day of the month
func MonthlyOnDay(daysOfMonth []int, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:       "monthly",
		DayOfMonth: daysOfMonth,
		TimeOfDay:  timeOfDay,
	}
}

func (r *Recurrence) IsWeekly() bool {
	return r.Type == "weekly" && len(r.Weekdays) > 0
}

func (r *Recurrence) IsInterval() bool {
	return r.Type == "interval" && r.Interval > 0
}

func (r *Recurrence) IsMonthly() bool {
	return r.Type == "monthly" && len(r.DayOfMonth) > 0
}

func (r *Recurrence) IsDaily() bool {
	return r.Type == "daily"
}
