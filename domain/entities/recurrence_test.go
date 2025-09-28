package entities

import (
	"testing"
	"time"
)

func TestParseTimeOfDay_ValidAndInvalid(t *testing.T) {
	// valid padded
	tm, err := parseTimeOfDay("09:05")
	if err != nil {
		t.Fatalf("expected parse success, got err: %v", err)
	}
	if tm.Hour() != 9 || tm.Minute() != 5 {
		t.Fatalf("expected 09:05, got %02d:%02d", tm.Hour(), tm.Minute())
	}

	// valid non-padded
	tm, err = parseTimeOfDay("9:05")
	if err != nil {
		t.Fatalf("expected parse success, got err: %v", err)
	}
	if tm.Hour() != 9 || tm.Minute() != 5 {
		t.Fatalf("expected 09:05, got %02d:%02d", tm.Hour(), tm.Minute())
	}

	// invalid
	_, err = parseTimeOfDay("bad")
	if err == nil {
		t.Fatalf("expected parse error for invalid string")
	}
}

func TestOnceAt_GetTimeOfDayAndLocation(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	rec := OnceAt(date, "07:30", loc)
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
	recW := CustomWeekly([]time.Weekday{time.Monday, time.Wednesday}, "08:20", loc)
	if !recW.IsWeekly() {
		t.Fatalf("expected weekly recurrence")
	}
	if recW.GetTimeOfDay() != "08:20" {
		t.Fatalf("expected weekly time 08:20, got %s", recW.GetTimeOfDay())
	}

	// Daily
	recD := DailyAt("14:45", loc)
	if !recD.IsDaily() {
		t.Fatalf("expected daily recurrence")
	}

	// Monthly
	recM := MonthlyOnDay([]int{1, 15}, "06:00", loc)
	if !recM.IsMonthly() {
		t.Fatalf("expected monthly recurrence")
	}
	if len(recM.DayOfMonth) != 2 || recM.DayOfMonth[0] != 1 || recM.DayOfMonth[1] != 15 {
		t.Fatalf("unexpected days of month: %v", recM.DayOfMonth)
	}

	// Interval
	recI := IntervalEveryDays(3, "05:05", loc)
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
