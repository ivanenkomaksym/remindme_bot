package scheduler

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

func mustTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	tm, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		t.Fatalf("parse time: %v", err)
	}
	return tm
}

func TestNextDailyTrigger(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	// Time later today -> same day
	got := NextDailyTrigger(base, "23:00")
	want := mustTime(t, "2006-01-02 15:04", "2025-01-10 23:00")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	// Time earlier today -> next day
	got = NextDailyTrigger(base, "09:45")
	want = mustTime(t, "2006-01-02 15:04", "2025-01-11 09:45")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextWeeklyTrigger(t *testing.T) {
	// Friday 2025-01-10 10:30
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	// Ask for Monday and Friday at 11:00 -> should pick same Friday 11:00
	got := NextWeeklyTrigger(base, []time.Weekday{time.Monday, time.Friday}, "11:00")
	want := mustTime(t, "2006-01-02 15:04", "2025-01-10 11:00")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	// Ask for Monday at 09:00 -> next Monday
	got = NextWeeklyTrigger(base, []time.Weekday{time.Monday}, "09:00")
	want = mustTime(t, "2006-01-02 15:04", "2025-01-13 09:00")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextMonthlyTrigger(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	// Ask for day 15 09:00 -> same month
	got := NextMonthlyTrigger(base, []int{15}, "09:00")
	want := mustTime(t, "2006-01-02 15:04", "2025-01-15 09:00")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	// Ask for day 1 09:00 -> next month
	got = NextMonthlyTrigger(base, []int{1}, "09:00")
	want = mustTime(t, "2006-01-02 15:04", "2025-02-01 09:00")
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNextForRecurrence(t *testing.T) {
	base := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30")
	// Daily advances 24h
	recDaily := models.DailyAt("10:30")
	got := NextForRecurrence(base, recDaily)
	want := base.Add(24 * time.Hour)
	if !got.Equal(want) {
		t.Fatalf("daily: expected %v, got %v", want, got)
	}
	// Weekly Monday at 09:00 from Friday -> next Monday at 09:00
	recWeekly := models.CustomWeekly([]time.Weekday{time.Monday}, "09:00")
	got = NextForRecurrence(base, recWeekly)
	want = mustTime(t, "2006-01-02 15:04", "2025-01-13 09:00")
	if !got.Equal(want) {
		t.Fatalf("weekly: expected %v, got %v", want, got)
	}
	// Monthly 15th 09:00
	recMonthly := models.MonthlyOnDay([]int{15}, "09:00")
	got = NextForRecurrence(base, recMonthly)
	want = mustTime(t, "2006-01-02 15:04", "2025-01-15 09:00")
	if !got.Equal(want) {
		t.Fatalf("monthly: expected %v, got %v", want, got)
	}
}
