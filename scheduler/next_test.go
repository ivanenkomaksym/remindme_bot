package scheduler

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func mustTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	// Parse in UTC to have deterministic base instants; tests convert to loc when needed.
	tm, err := time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		t.Fatalf("parse time: %v", err)
	}
	return tm
}

func TestNextDailyTrigger(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// Time later today -> same day
	var wantLocal time.Time
	var want time.Time
	got := NextDailyTrigger(base, "23:00", loc)
	wantLocal = time.Date(2025, 1, 10, 23, 0, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	// Time earlier today -> next day
	got = NextDailyTrigger(base, "09:45", loc)
	wantLocal = time.Date(2025, 1, 11, 9, 45, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextWeeklyTrigger(t *testing.T) {
	// Friday 2025-01-10 10:30 (as local time)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// Build base as the local instant 2025-01-10 10:30 in user's tz, convert to UTC for scheduler input
	baseLocal := time.Date(2025, 1, 10, 10, 30, 0, 0, loc)
	base := baseLocal.UTC()

	// Ask for Monday and Friday at 11:00 -> should pick same Friday 11:00
	got := NextWeeklyTrigger(base, []time.Weekday{time.Monday, time.Friday}, "11:00", loc)
	wantLocal := time.Date(2025, 1, 10, 11, 0, 0, 0, loc)
	want := wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	// Ask for Monday at 09:00 -> next Monday
	got = NextWeeklyTrigger(base, []time.Weekday{time.Monday}, "09:00", loc)
	wantLocal = time.Date(2025, 1, 13, 9, 0, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextMonthlyTrigger(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// Ask for day 15 09:00 -> same month
	got := NextMonthlyTrigger(base, []int{15}, "09:00", loc)
	wantLocal := time.Date(2025, 1, 15, 9, 0, 0, 0, loc)
	want := wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	// Ask for day 1 09:00 -> next month
	got = NextMonthlyTrigger(base, []int{1}, "09:00", loc)
	wantLocal = time.Date(2025, 2, 1, 9, 0, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextForRecurrence(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	var wantLocal time.Time
	var want time.Time

	recOnce := entities.OnceAt(base, "10:30", loc)
	got := NextForRecurrence(base, recOnce)
	if got != nil {
		t.Fatalf("daily: expected %v, got %v", nil, got)
	}

	// Daily advances to the next local 10:30 -> compute expected candidate in loc and compare in UTC
	recDaily := entities.DailyAt("10:30", loc)
	got = NextForRecurrence(base, recDaily)
	baseInLoc := base.In(loc)
	candidate := time.Date(baseInLoc.Year(), baseInLoc.Month(), baseInLoc.Day(), 10, 30, 0, 0, loc)
	if !candidate.After(baseInLoc) {
		candidate = candidate.Add(24 * time.Hour)
	}
	want = candidate.UTC()
	if !got.Equal(want) {
		t.Fatalf("daily: expected %v, got %v", want, got)
	}
	// Weekly Monday at 09:00 from Friday -> next Monday at 09:00
	recWeekly := entities.CustomWeekly([]time.Weekday{time.Monday}, "09:00", loc)
	got = NextForRecurrence(base, recWeekly)
	wantLocal = time.Date(2025, 1, 13, 9, 0, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("weekly: expected %v, got %v", want, got)
	}
	// Monthly 15th 09:00
	recMonthly := entities.MonthlyOnDay([]int{15}, "09:00", loc)
	got = NextForRecurrence(base, recMonthly)
	wantLocal = time.Date(2025, 1, 15, 9, 0, 0, 0, loc)
	want = wantLocal.UTC()
	if !got.Equal(want) {
		t.Fatalf("monthly: expected %v, got %v", want, got)
	}
}
