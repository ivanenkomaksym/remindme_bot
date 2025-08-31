package models

import (
	"testing"
	"time"
)

func TestRecurrenceExamples(t *testing.T) {
	// Example 1: Gym on Tuesday and Thursday at 6:00 AM
	gymRecurrence := GymOnTueThu("06:00")
	if gymRecurrence.Type != Weekly {
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
	if classesRecurrence.Type != Weekly {
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
	if medicineRecurrence.Type != Interval {
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
	if cleaningRecurrence.Type != Weekly {
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

	// Example 6: Monthly on the 1st and 15th day of the month
	monthlyRecurrence := &Recurrence{Type: Monthly, DayOfMonth: []int{1, 15}}
	if !monthlyRecurrence.IsMonthly() {
		t.Error("Expected IsMonthly() to return true")
	}
	if len(monthlyRecurrence.DayOfMonth) != 2 {
		t.Errorf("Expected 2 days of month, got %d", len(monthlyRecurrence.DayOfMonth))
	}
	if monthlyRecurrence.DayOfMonth[0] != 1 || monthlyRecurrence.DayOfMonth[1] != 15 {
		t.Errorf("Expected 1 and 15, got %v", monthlyRecurrence.DayOfMonth)
	}
}

func TestRecurrenceHelperMethods(t *testing.T) {
	// Test helper methods
	weeklyRecurrence := &Recurrence{Type: Weekly, Weekdays: []time.Weekday{time.Monday}}
	if !weeklyRecurrence.IsWeekly() {
		t.Error("Expected IsWeekly() to return true")
	}

	intervalRecurrence := &Recurrence{Type: Interval, Interval: 3}
	if !intervalRecurrence.IsInterval() {
		t.Error("Expected IsInterval() to return true")
	}

	dailyRecurrence := &Recurrence{Type: Daily}
	if !dailyRecurrence.IsDaily() {
		t.Error("Expected IsDaily() to return true")
	}

	monthlyRecurrence := &Recurrence{Type: Monthly, DayOfMonth: []int{1, 15}}
	if !monthlyRecurrence.IsMonthly() {
		t.Error("Expected IsMonthly() to return true")
	}
}

// GymOnTueThu creates a recurrence for gym on Tuesday and Thursday
func GymOnTueThu(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Weekly,
		Weekdays:  []time.Weekday{time.Tuesday, time.Thursday},
		TimeOfDay: timeOfDay,
	}
}

// ClassesOnWedFri creates a recurrence for classes on Wednesday and Friday
func ClassesOnWedFri(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Weekly,
		Weekdays:  []time.Weekday{time.Wednesday, time.Friday},
		TimeOfDay: timeOfDay,
	}
}

// MedicineEveryNDays creates a recurrence for medicine every N days
func MedicineEveryNDays(interval int, timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Interval,
		Interval:  interval,
		TimeOfDay: timeOfDay,
	}
}

// CleaningEveryFriday creates a recurrence for cleaning every Friday
func CleaningEveryFriday(timeOfDay string) *Recurrence {
	return &Recurrence{
		Type:      Weekly,
		Weekdays:  []time.Weekday{time.Friday},
		TimeOfDay: timeOfDay,
	}
}
