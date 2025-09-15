package entities

import "time"

type Recurrence struct {
	Type       RecurrenceType `json:"type"`          // e.g., "daily", "weekly", "monthly", "interval", "custom"
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
		Type:      Weekly,
		Weekdays:  weekdays,
		TimeOfDay: timeOfDay,
	}
}

func OnceAt(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Once,
		TimeOfDay: timeOfDay,
	}
}

// DailyAt creates a daily recurrence at a specific time
func DailyAt(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Daily,
		TimeOfDay: timeOfDay,
	}
}

// MonthlyOnDay creates a monthly recurrence on a specific day of the month
func MonthlyOnDay(daysOfMonth []int, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:       Monthly,
		DayOfMonth: daysOfMonth,
		TimeOfDay:  timeOfDay,
	}
}

func (r *Recurrence) IsWeekly() bool {
	return r.Type == Weekly && len(r.Weekdays) > 0
}

func (r *Recurrence) IsInterval() bool {
	return r.Type == Interval && r.Interval > 0
}

func (r *Recurrence) IsMonthly() bool {
	return r.Type == Monthly && len(r.DayOfMonth) > 0
}

func (r *Recurrence) IsDaily() bool {
	return r.Type == Daily
}
