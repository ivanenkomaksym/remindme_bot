package models

import (
	"testing"
	"time"
)

func TestRecurrenceExamples(t *testing.T) {
	// Example 1: Gym on Tuesday and Thursday at 6:00 AM
	gymRecurrence := GymOnTueThu("06:00")
	if gymRecurrence.Type != "weekly" {
		t.Errorf("Expected weekly type, got %s", gymRecurrence.Type)
	}
	if len(gymRecurrence.Weekdays) != 2 {
		t.Errorf("Expected 2 weekdays, got %d", len(gymRecurrence.Weekdays))
	}
	if gymRecurrence.Weekdays[0] != time.Tuesday || gymRecurrence.Weekdays[1] != time.Thursday {
		t.Errorf("Expected Tuesday and Thursday, got %v", gymRecurrence.Weekdays)
	}
	if gymRecurrence.TimeOfDay != "06:00" {
		t.Errorf("Expected 06:00, got %s", gymRecurrence.TimeOfDay)
	}

	// Example 2: Classes on Wednesday and Friday at 2:00 PM
	classesRecurrence := ClassesOnWedFri("14:00")
	if classesRecurrence.Type != "weekly" {
		t.Errorf("Expected weekly type, got %s", classesRecurrence.Type)
	}
	if len(classesRecurrence.Weekdays) != 2 {
		t.Errorf("Expected 2 weekdays, got %d", len(classesRecurrence.Weekdays))
	}
	if classesRecurrence.Weekdays[0] != time.Wednesday || classesRecurrence.Weekdays[1] != time.Friday {
		t.Errorf("Expected Wednesday and Friday, got %v", classesRecurrence.Weekdays)
	}

	// Example 3: Medicine every 5 days at 9:00 AM
	medicineRecurrence := MedicineEveryNDays(5, "09:00")
	if medicineRecurrence.Type != "interval" {
		t.Errorf("Expected interval type, got %s", medicineRecurrence.Type)
	}
	if medicineRecurrence.Interval != 5 {
		t.Errorf("Expected interval 5, got %d", medicineRecurrence.Interval)
	}
	if medicineRecurrence.TimeOfDay != "09:00" {
		t.Errorf("Expected 09:00, got %s", medicineRecurrence.TimeOfDay)
	}

	// Example 4: Cleaning every Friday at 10:00 AM
	cleaningRecurrence := CleaningEveryFriday("10:00")
	if cleaningRecurrence.Type != "weekly" {
		t.Errorf("Expected weekly type, got %s", cleaningRecurrence.Type)
	}
	if len(cleaningRecurrence.Weekdays) != 1 {
		t.Errorf("Expected 1 weekday, got %d", len(cleaningRecurrence.Weekdays))
	}
	if cleaningRecurrence.Weekdays[0] != time.Friday {
		t.Errorf("Expected Friday, got %v", cleaningRecurrence.Weekdays[0])
	}

	// Example 5: Custom weekly pattern with start and end dates
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	customRecurrence := CustomWeekly([]time.Weekday{time.Monday, time.Wednesday, time.Friday}, "15:00").
		WithStartDate(startDate).
		WithEndDate(endDate)

	if customRecurrence.StartDate == nil || customRecurrence.EndDate == nil {
		t.Error("Expected start and end dates to be set")
	}
	if len(customRecurrence.Weekdays) != 3 {
		t.Errorf("Expected 3 weekdays, got %d", len(customRecurrence.Weekdays))
	}
}

func TestRecurrenceHelperMethods(t *testing.T) {
	// Test helper methods
	weeklyRecurrence := &Recurrence{Type: "weekly", Weekdays: []time.Weekday{time.Monday}}
	if !weeklyRecurrence.IsWeekly() {
		t.Error("Expected IsWeekly() to return true")
	}

	intervalRecurrence := &Recurrence{Type: "interval", Interval: 3}
	if !intervalRecurrence.IsInterval() {
		t.Error("Expected IsInterval() to return true")
	}

	dailyRecurrence := &Recurrence{Type: "daily"}
	if !dailyRecurrence.IsDaily() {
		t.Error("Expected IsDaily() to return true")
	}

	monthlyRecurrence := &Recurrence{Type: "monthly", DayOfMonth: 15}
	if !monthlyRecurrence.IsMonthly() {
		t.Error("Expected IsMonthly() to return true")
	}
}

// GymOnTueThu creates a recurrence for gym on Tuesday and Thursday
func GymOnTueThu(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "weekly",
		Weekdays:  []time.Weekday{time.Tuesday, time.Thursday},
		TimeOfDay: timeOfDay,
	}
}

// ClassesOnWedFri creates a recurrence for classes on Wednesday and Friday
func ClassesOnWedFri(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "weekly",
		Weekdays:  []time.Weekday{time.Wednesday, time.Friday},
		TimeOfDay: timeOfDay,
	}
}

// MedicineEveryNDays creates a recurrence for medicine every N days
func MedicineEveryNDays(interval int, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "interval",
		Interval:  interval,
		TimeOfDay: timeOfDay,
	}
}

// CleaningEveryFriday creates a recurrence for cleaning every Friday
func CleaningEveryFriday(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      "weekly",
		Weekdays:  []time.Weekday{time.Friday},
		TimeOfDay: timeOfDay,
	}
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
func MonthlyOnDay(dayOfMonth int, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:       "monthly",
		DayOfMonth: dayOfMonth,
		TimeOfDay:  timeOfDay,
	}
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

// WithOccurrences adds a limit on the number of occurrences
func (r *Recurrence) WithOccurrences(occurrences int) *Recurrence {
	r.Occurrences = &occurrences
	return r
}
