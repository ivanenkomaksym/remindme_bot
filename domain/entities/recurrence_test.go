package entities

import (
	"testing"
	"time"
)

func TestOnceAt_GetTimeOfDayAndLocation(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(2025, 1, 10, 7, 30, 0, 0, loc)
	rec := OnceAt(date, loc)
	if rec == nil {
		t.Fatalf("expected recurrence")
	}
	if rec.Type != Once {
		t.Fatalf("expected Type Once, got %v", rec.Type)
	}
	got := rec.GetTimeOfDay()
	if got != "07:30" {
		t.Fatalf("expected GetTimeOfDay 07:30, got %s", got)
	}
	// Location should be set
	if rec.GetLocation() == nil {
		t.Fatalf("expected location to be set")
	}
}

func TestCustomWeeklyAndDailyAndMonthlyAndInterval(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// Custom weekly
	timeOfDay := time.Date(2025, 1, 10, 8, 20, 0, 0, loc)
	recW := CustomWeekly([]time.Weekday{time.Monday, time.Wednesday}, timeOfDay, loc)
	if !recW.IsWeekly() {
		t.Fatalf("expected weekly recurrence")
	}
	if recW.GetTimeOfDay() != "08:20" {
		t.Fatalf("expected weekly time 08:20, got %s", recW.GetTimeOfDay())
	}

	// Daily
	timeOfDay = time.Date(2025, 1, 10, 14, 45, 0, 0, loc)
	recD := DailyAt(timeOfDay, loc)
	if !recD.IsDaily() {
		t.Fatalf("expected daily recurrence")
	}

	// Monthly
	timeOfDay = time.Date(2025, 1, 10, 6, 0, 0, 0, loc)
	recM := MonthlyOnDay([]int{1, 15}, timeOfDay, loc)
	if !recM.IsMonthly() {
		t.Fatalf("expected monthly recurrence")
	}
	if len(recM.DayOfMonth) != 2 || recM.DayOfMonth[0] != 1 || recM.DayOfMonth[1] != 15 {
		t.Fatalf("unexpected days of month: %v", recM.DayOfMonth)
	}

	// Interval
	timeOfDay = time.Date(2025, 1, 10, 5, 5, 0, 0, loc)
	recI := IntervalEveryDays(3, timeOfDay, loc)
	if !recI.IsInterval() {
		t.Fatalf("expected interval recurrence")
	}
	if recI.Interval != 3 {
		t.Fatalf("expected interval 3, got %d", recI.Interval)
	}
}

func TestGetLocation_LoadsFromName(t *testing.T) {
	// Prepare a recurrence with LocationName set but Location nil
	rec := &Recurrence{LocationName: "Asia/Shanghai"}
	loc := rec.GetLocation()
	if loc == nil {
		t.Fatalf("expected location to be loaded from name")
	}
}
