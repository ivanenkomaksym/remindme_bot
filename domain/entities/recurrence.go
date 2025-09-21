package entities

import "time"

type Recurrence struct {
	Type       RecurrenceType `json:"type"`          // e.g., "once", "daily", "weekly", "monthly", "interval", "custom"
	Interval   int            `json:"interval"`      // e.g., every N days/hours/minutes
	Weekdays   []time.Weekday `json:"weekdays"`      // For weekly recurrence (e.g., [Tuesday, Thursday])
	DayOfMonth []int          `json:"days_of_month"` // For monthly recurrence (e.g., [1, 15])
	StartDate  *time.Time     `json:"start_date"`    // When recurrence begins (includes time)
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

// GetTimeOfDay returns the time of day in "HH:MM" format from StartDate
func (r *Recurrence) GetTimeOfDay() string {
	if r.StartDate == nil {
		return "00:00"
	}
	return r.StartDate.Format("15:04")
}

// GetTimeOfDayAsTime returns the time portion of StartDate
func (r *Recurrence) GetTimeOfDayAsTime() time.Time {
	if r.StartDate == nil {
		return time.Time{}
	}
	hour, minute, second := r.StartDate.Clock()
	return time.Date(0, 1, 1, hour, minute, second, 0, time.UTC)
}

// CustomWeekly creates a custom weekly recurrence on specific weekdays
func CustomWeekly(weekdays []time.Weekday, timeOfDay string) *Recurrence {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Parse time and set it to startDate
	if hour, minute, ok := parseTimeOfDay(timeOfDay); ok {
		startDate = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	}

	return &Recurrence{
		Type:      Weekly,
		Weekdays:  weekdays,
		StartDate: &startDate,
	}
}

func OnceAt(date time.Time, timeOfDay string) *Recurrence {
	// Parse time and set it to the provided date
	startDate := date
	if hour, minute, ok := parseTimeOfDay(timeOfDay); ok {
		startDate = time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
	}

	return &Recurrence{
		Type:      Once,
		StartDate: &startDate,
	}
}

// DailyAt creates a daily recurrence at a specific time
func DailyAt(timeOfDay string) *Recurrence {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Parse time and set it to startDate
	if hour, minute, ok := parseTimeOfDay(timeOfDay); ok {
		startDate = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	}

	return &Recurrence{
		Type:      Daily,
		StartDate: &startDate,
	}
}

// MonthlyOnDay creates a monthly recurrence on a specific day of the month
func MonthlyOnDay(daysOfMonth []int, timeOfDay string) *Recurrence {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Parse time and set it to startDate
	if hour, minute, ok := parseTimeOfDay(timeOfDay); ok {
		startDate = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	}

	return &Recurrence{
		Type:       Monthly,
		DayOfMonth: daysOfMonth,
		StartDate:  &startDate,
	}
}

// parseTimeOfDay parses a time string in "HH:MM" format
func parseTimeOfDay(timeStr string) (hour, minute int, ok bool) {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return 0, 0, false
	}

	hourStr := timeStr[:2]
	minuteStr := timeStr[3:]

	if h, err := time.Parse("15", hourStr); err != nil {
		return 0, 0, false
	} else if m, err := time.Parse("04", minuteStr); err != nil {
		return 0, 0, false
	} else {
		return h.Hour(), m.Minute(), true
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
