package entities

import (
	"time"
)

type Recurrence struct {
	Type                      RecurrenceType `json:"type" bson:"type"`                   // e.g., "once", "daily", "weekly", "monthly", "interval", "custom"
	Interval                  int            `json:"interval" bson:"interval"`           // e.g., every N days/hours/minutes
	Weekdays                  []time.Weekday `json:"weekdays" bson:"weekdays"`           // For weekly recurrence (e.g., [Tuesday, Thursday])
	DayOfMonth                []int          `json:"days_of_month" bson:"days_of_month"` // For monthly recurrence (e.g., [1, 15])
	StartDate                 *time.Time     `json:"start_date" bson:"start_date"`       // When recurrence begins (includes time)
	LocationName              string         `json:"location" bson:"location"`
	Location                  *time.Location `json:"-" bson:"-"`                                                       // Ignore
	EndDate                   *time.Time     `json:"end_date" bson:"end_date"`                                         // When recurrence ends (optional)
	SpacedBasedRepetitionDays []int          `json:"spaced_based_repetition_days" bson:"spaced_based_repetition_days"` // For spaced-based repetition (e.g., [1, 3, 7, 14])
}

type Option func(dp *Recurrence)

func WithInterval(interval int) Option {
	return func(r *Recurrence) {
		r.Interval = interval
	}
}

func WithWeekdays(weekdays []time.Weekday) Option {
	return func(r *Recurrence) {
		r.Weekdays = weekdays
	}
}

func WithDaysOfMonth(daysOfMonth []int) Option {
	return func(r *Recurrence) {
		r.DayOfMonth = daysOfMonth
	}
}

func WithSpacedBasedRepetition() Option {
	return func(r *Recurrence) {
		r.SpacedBasedRepetitionDays = []int{0, 1, 2, 3, 5, 7, 7, 7}
	}
}

func New(recurrenceType RecurrenceType, startDate *time.Time, location *time.Location, opts ...Option) *Recurrence {
	recurrence := &Recurrence{
		Type:      recurrenceType,
		StartDate: startDate,
	}

	recurrence.SetLocation(location)

	for _, opt := range opts {
		opt(recurrence)
	}

	return recurrence
}

func (r *Recurrence) GetLocation() *time.Location {
	// If the private field is nil, try to load it from the stored string.
	if r.Location == nil && r.LocationName != "" {
		loc, err := time.LoadLocation(r.LocationName)
		if err == nil {
			r.Location = loc
		}
	}
	return r.Location
}

func (r *Recurrence) SetLocation(loc *time.Location) {
	r.Location = loc
	if loc != nil {
		r.LocationName = loc.String()
	} else {
		r.LocationName = ""
	}
}

// GetTimeOfDay returns the time of day in "HH:MM" format from StartDate
func (r *Recurrence) GetTimeOfDay() string {
	if r.StartDate == nil {
		return "00:00"
	}
	return r.StartDate.In(r.GetLocation()).Format("15:04")
}

// CustomWeekly creates a custom weekly recurrence on specific weekdays
func CustomWeekly(weekdays []time.Weekday, timeOfDay time.Time, location *time.Location) *Recurrence {
	return New(Weekly, &timeOfDay, location, WithWeekdays(weekdays))
}

func OnceAt(startDate time.Time, location *time.Location) *Recurrence {
	return New(Once, &startDate, location)
}

// DailyAt creates a daily recurrence at a specific time
func DailyAt(timeOfDay time.Time, location *time.Location) *Recurrence {
	return New(Daily, &timeOfDay, location)
}

// IntervalEveryDays creates a recurrence that triggers every N days at a specific time
func IntervalEveryDays(intervalDays int, timeOfDay time.Time, location *time.Location) *Recurrence {
	return New(Interval, &timeOfDay, location, WithInterval(intervalDays))
}

// MonthlyOnDay creates a monthly recurrence on a specific day of the month
func MonthlyOnDay(daysOfMonth []int, timeOfDay time.Time, location *time.Location) *Recurrence {
	return New(Monthly, &timeOfDay, location, WithDaysOfMonth(daysOfMonth))
}

func SpacedBasedRepetitionInterval(timeOfDay time.Time, location *time.Location) *Recurrence {
	return New(SpacedBasedRepetition, &timeOfDay, location, WithSpacedBasedRepetition())
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
