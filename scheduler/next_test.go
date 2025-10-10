package scheduler

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func mustTime(t *testing.T, layout, value string, location *time.Location) time.Time {
	t.Helper()
	tm, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		t.Fatalf("parse time: %v", err)
	}
	return tm
}

func TestNextDailyTrigger(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	from := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30", loc)

	// Time later today -> same day
	wantLocal := time.Date(2025, 1, 10, 23, 0, 0, 0, loc)
	wantUTC := wantLocal.UTC()

	timeOfDay := time.Date(2025, 1, 10, 23, 0, 0, 0, loc)
	gotUTC := NextDailyTrigger(from, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}

	// Time earlier today -> next day
	wantLocal = time.Date(2025, 1, 11, 9, 45, 0, 0, loc)
	wantUTC = wantLocal.UTC()

	timeOfDay = time.Date(2025, 1, 10, 9, 45, 0, 0, loc)
	gotUTC = NextDailyTrigger(from, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}
}

func TestNextWeeklyTrigger(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// Build base as the local instant 2025-01-10 10:30 in user's tz
	from := time.Date(2025, 1, 10, 10, 30, 0, 0, loc)

	// Ask for Monday and Friday at 11:00 -> should pick same Friday 11:00
	wantLocal := time.Date(2025, 1, 10, 11, 0, 0, 0, loc)
	wantUTC := wantLocal.UTC()

	timeOfDay := time.Date(2025, 1, 10, 11, 00, 0, 0, loc)
	gotUTC := NextWeeklyTrigger(from, []time.Weekday{time.Monday, time.Friday}, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}
	// Ask for Monday at 09:00 -> next Monday
	wantLocal = time.Date(2025, 1, 13, 9, 0, 0, 0, loc)
	wantUTC = wantLocal.UTC()

	timeOfDay = time.Date(2025, 1, 10, 9, 00, 0, 0, loc)
	gotUTC = NextWeeklyTrigger(from, []time.Weekday{time.Monday}, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}
}

func TestNextMonthlyTrigger(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	from := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30", loc)

	// Ask for day 15 09:00 -> same month
	wantLocal := time.Date(2025, 1, 15, 9, 0, 0, 0, loc)
	wantUTC := wantLocal.UTC()

	timeOfDay := time.Date(2025, 1, 10, 9, 00, 0, 0, loc)
	gotUTC := NextMonthlyTrigger(from, []int{15}, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}

	// Ask for day 1 09:00 -> next month
	wantLocal = time.Date(2025, 2, 1, 9, 0, 0, 0, loc)
	wantUTC = wantLocal.UTC()

	gotUTC = NextMonthlyTrigger(from, []int{1}, timeOfDay, loc)
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("expected %v, got %v", wantUTC, gotUTC)
	}
}

func TestNextForRecurrence(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	from := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30", loc)

	var wantLocal time.Time
	var wantUTC time.Time

	timeOfDay := time.Date(2025, 1, 10, 10, 30, 0, 0, loc)
	recOnce := entities.OnceAt(timeOfDay, loc)
	gotUTC := NextForRecurrence(from, timeOfDay, recOnce)
	if gotUTC != nil {
		t.Fatalf("daily: expected %v, got %v", nil, gotUTC)
	}

	// Daily advances to the next local 10:30 -> compute expected candidate in loc and compare in UTC
	recDaily := entities.DailyAt(timeOfDay, loc)
	gotUTC = NextForRecurrence(from, timeOfDay, recDaily)
	candidate := time.Date(from.Year(), from.Month(), from.Day(), 10, 30, 0, 0, loc)
	if !candidate.After(from) {
		candidate = candidate.Add(24 * time.Hour)
	}
	wantUTC = candidate.UTC()
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("daily: expected %v, got %v", wantUTC, gotUTC)
	}
	// Weekly Monday at 09:00 from Friday -> next Monday at 09:00
	timeOfDay = time.Date(2025, 1, 10, 9, 0, 0, 0, loc)
	recWeekly := entities.CustomWeekly([]time.Weekday{time.Monday}, timeOfDay, loc)
	gotUTC = NextForRecurrence(from, timeOfDay, recWeekly)

	wantLocal = time.Date(2025, 1, 13, 9, 0, 0, 0, loc)
	wantUTC = wantLocal.UTC()
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("weekly: expected %v, got %v", wantUTC, gotUTC)
	}
	// Monthly 15th 09:00
	recMonthly := entities.MonthlyOnDay([]int{15}, timeOfDay, loc)
	gotUTC = NextForRecurrence(from, timeOfDay, recMonthly)
	wantLocal = time.Date(2025, 1, 15, 9, 0, 0, 0, loc)
	wantUTC = wantLocal.UTC()
	if !gotUTC.Equal(wantUTC) {
		t.Fatalf("monthly: expected %v, got %v", wantUTC, gotUTC)
	}
}
