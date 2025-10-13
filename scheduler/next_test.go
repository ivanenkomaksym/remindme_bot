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
	loc, _ := time.LoadLocation("Europe/Kyiv")
	// Friday January 10, 2025 at 10:30 AM Kyiv time (avoid edge cases)
	from := mustTime(t, "2006-01-02 15:04", "2025-01-10 10:30", loc)

	t.Run("Daily recurrence same day", func(t *testing.T) {
		// Want daily at 11:00 AM Kyiv time (same day since 11:00 > 10:30)
		timeOfDay := time.Date(2025, 1, 10, 11, 0, 0, 0, loc)
		dailyRec := entities.DailyAt(timeOfDay, loc)
		
		gotUTC := NextForRecurrence(from, timeOfDay, dailyRec)
		if gotUTC == nil {
			t.Fatal("Expected non-nil result for daily recurrence")
		}
		
		// Should be today at 11:00 AM Kyiv time
		wantLocal := time.Date(2025, 1, 10, 11, 0, 0, 0, loc)
		wantUTC := wantLocal.UTC()
		
		if !gotUTC.Equal(wantUTC) {
			t.Errorf("Daily same day: expected %v, got %v", wantUTC, *gotUTC)
		}
	})

	t.Run("Daily recurrence next day", func(t *testing.T) {
		// Want daily at 9:00 AM Kyiv time (next day since 9:00 < 10:30)
		timeOfDay := time.Date(2025, 1, 10, 9, 0, 0, 0, loc)
		dailyRec := entities.DailyAt(timeOfDay, loc)
		
		gotUTC := NextForRecurrence(from, timeOfDay, dailyRec)
		if gotUTC == nil {
			t.Fatal("Expected non-nil result for daily recurrence")
		}
		
		// Should be tomorrow at 9:00 AM Kyiv time
		wantLocal := time.Date(2025, 1, 11, 9, 0, 0, 0, loc)
		wantUTC := wantLocal.UTC()
		
		if !gotUTC.Equal(wantUTC) {
			t.Errorf("Daily next day: expected %v, got %v", wantUTC, *gotUTC)
		}
	})

	t.Run("Monthly recurrence same month", func(t *testing.T) {
		// Want 15th of month at 9:00 AM (same month since 15 > 10)
		timeOfDay := time.Date(2025, 1, 10, 9, 0, 0, 0, loc)
		recMonthly := entities.MonthlyOnDay([]int{15}, timeOfDay, loc)
		
		gotUTC := NextForRecurrence(from, timeOfDay, recMonthly)
		if gotUTC == nil {
			t.Fatal("Expected non-nil result for monthly recurrence")
		}
		
		// Should be January 15, 2025 at 9:00 AM Kyiv time
		wantLocal := time.Date(2025, 1, 15, 9, 0, 0, 0, loc)
		wantUTC := wantLocal.UTC()
		
		if !gotUTC.Equal(wantUTC) {
			t.Errorf("Monthly same month: expected %v, got %v", wantUTC, *gotUTC)
		}
	})

	t.Run("Interval recurrence", func(t *testing.T) {
		timeOfDay := time.Date(2025, 1, 10, 9, 0, 0, 0, loc)
		recInterval := &entities.Recurrence{
			Type:     entities.Interval,
			Interval: 3, // Every 3 days
		}
		recInterval.SetLocation(loc)
		
		gotUTC := NextForRecurrence(from, timeOfDay, recInterval)
		if gotUTC == nil {
			t.Fatal("Expected non-nil result for interval recurrence")
		}
		
		// Should advance by 3 days from the from time
		fromInTz := from.In(loc)
		wantLocal := fromInTz.Add(3 * 24 * time.Hour)
		wantUTC := wantLocal.UTC()
		
		if !gotUTC.Equal(wantUTC) {
			t.Errorf("Interval: expected %v, got %v", wantUTC, *gotUTC)
		}
	})

	t.Run("Once recurrence returns nil", func(t *testing.T) {
		timeOfDay := time.Date(2025, 1, 10, 9, 0, 0, 0, loc)
		recOnce := &entities.Recurrence{Type: entities.Once}
		
		gotUTC := NextForRecurrence(from, timeOfDay, recOnce)
		if gotUTC != nil {
			t.Errorf("Once: expected nil, got %v", *gotUTC)
		}
	})
}