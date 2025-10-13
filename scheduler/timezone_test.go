package scheduler

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// TestCrossTimezoneScheduling tests that scheduling works correctly when
// the input times are in different timezones than the target location
func TestCrossTimezoneScheduling(t *testing.T) {
	// User in New York timezone
	userTZ, _ := time.LoadLocation("America/New_York")

	// From time in UTC
	fromUTC := time.Date(2025, 1, 10, 15, 30, 0, 0, time.UTC) // 3:30 PM UTC

	// Time of day also in UTC (should be converted to user's timezone)
	timeOfDayUTC := time.Date(2025, 1, 10, 18, 0, 0, 0, time.UTC) // 6:00 PM UTC

	t.Run("NextDailyTrigger with timezone conversion", func(t *testing.T) {
		result := NextDailyTrigger(fromUTC, timeOfDayUTC, userTZ)

		// From: 3:30 PM UTC = 10:30 AM EST
		// TimeOfDay: 6:00 PM UTC = 1:00 PM EST
		// Since 1:00 PM EST is after 10:30 AM EST, it should be today
		expected := time.Date(2025, 1, 10, 18, 0, 0, 0, time.UTC)

		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("NextWeeklyTrigger with timezone conversion", func(t *testing.T) {
		// Friday in UTC, want next Monday at 6 PM UTC (1 PM EST)
		days := []time.Weekday{time.Monday}
		result := NextWeeklyTrigger(fromUTC, days, timeOfDayUTC, userTZ)

		// Should be Monday Jan 13, 2025 at 6 PM UTC (1 PM EST)
		expected := time.Date(2025, 1, 13, 18, 0, 0, 0, time.UTC)

		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("NextMonthlyTrigger with timezone conversion", func(t *testing.T) {
		daysOfMonth := []int{15}
		result := NextMonthlyTrigger(fromUTC, daysOfMonth, timeOfDayUTC, userTZ)

		// Should be Jan 15, 2025 at 6 PM UTC (1 PM EST)
		expected := time.Date(2025, 1, 15, 18, 0, 0, 0, time.UTC)

		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestDSTTransitions tests scheduling around Daylight Saving Time transitions
func TestDSTTransitions(t *testing.T) {
	// US Eastern timezone for DST testing
	eastern, _ := time.LoadLocation("America/New_York")

	// Spring forward: March 9, 2025 (2:00 AM becomes 3:00 AM)
	t.Run("Daily trigger across spring DST", func(t *testing.T) {
		// Saturday March 8, 2025 at 1:30 AM EST
		from := time.Date(2025, 3, 8, 1, 30, 0, 0, eastern)
		// Want 2:30 AM EST, but since we're at 1:30 AM, this should be today, not tomorrow
		timeOfDay := time.Date(2025, 3, 8, 2, 30, 0, 0, eastern)

		result := NextDailyTrigger(from, timeOfDay, eastern)

		// Since from (1:30 AM) is before timeOfDay (2:30 AM), should be same day
		resultInEastern := result.In(eastern)
		if resultInEastern.Hour() != 2 || resultInEastern.Minute() != 30 {
			t.Errorf("Expected 2:30 AM EST, got %v", resultInEastern.Format("15:04 MST"))
		}

		// Should be March 8th, not 9th
		if resultInEastern.Day() != 8 {
			t.Errorf("Expected March 8, got day %d", resultInEastern.Day())
		}
	})

	// Fall back: November 2, 2025 (2:00 AM becomes 1:00 AM)
	t.Run("Daily trigger across fall DST", func(t *testing.T) {
		// Saturday November 1, 2025 at 1:30 AM EDT
		from := time.Date(2025, 11, 1, 1, 30, 0, 0, eastern)
		// Want 1:30 AM (which exists twice on November 2)
		timeOfDay := time.Date(2025, 11, 1, 1, 30, 0, 0, eastern)

		result := NextDailyTrigger(from, timeOfDay, eastern)

		// Should advance to next day at 1:30 AM EST (after DST ends)
		resultInEastern := result.In(eastern)
		if resultInEastern.Hour() != 1 || resultInEastern.Minute() != 30 {
			t.Errorf("Expected 1:30 AM EST, got %v", resultInEastern.Format("15:04 MST"))
		}

		// Verify it's the next day
		if resultInEastern.Day() != 2 {
			t.Errorf("Expected November 2, got day %d", resultInEastern.Day())
		}
	})
}

// TestParseHourMinuteEdgeCases tests edge cases for time parsing
func TestParseHourMinuteEdgeCases(t *testing.T) {
	tests := []struct {
		input  string
		hour   int
		minute int
		valid  bool
		name   string
	}{
		{"00:00", 0, 0, true, "midnight"},
		{"23:59", 23, 59, true, "end of day"},
		{"12:30", 12, 30, true, "noon-ish"},
		{"24:00", 0, 0, false, "invalid hour 24"},
		{"-1:30", 0, 0, false, "negative hour"},
		{"12:60", 0, 0, false, "invalid minute 60"},
		{"12:-1", 0, 0, false, "negative minute"},
		{"12", 0, 0, false, "missing colon"},
		{"12:30:45", 0, 0, false, "seconds included"},
		{"", 0, 0, false, "empty string"},
		{"ab:cd", 0, 0, false, "non-numeric"},
		{"12:3a", 0, 0, false, "partial non-numeric"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, valid := ParseHourMinute(tt.input)
			if valid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v for input %q", tt.valid, valid, tt.input)
			}
			if valid && (hour != tt.hour || minute != tt.minute) {
				t.Errorf("Expected %d:%d, got %d:%d for input %q", tt.hour, tt.minute, hour, minute, tt.input)
			}
		})
	}
}

// TestNextForRecurrenceEdgeCases tests edge cases in the main scheduling function
func TestNextForRecurrenceEdgeCases(t *testing.T) {
	loc := time.UTC
	from := time.Date(2025, 1, 10, 12, 0, 0, 0, loc)
	timeOfDay := time.Date(2025, 1, 10, 15, 0, 0, 0, loc)

	t.Run("Once type returns nil", func(t *testing.T) {
		rec := &entities.Recurrence{Type: entities.Once}
		result := NextForRecurrence(from, timeOfDay, rec)
		if result != nil {
			t.Error("Expected nil for Once recurrence type")
		}
	})

	t.Run("Interval with zero days defaults to 1", func(t *testing.T) {
		rec := &entities.Recurrence{
			Type:     entities.Interval,
			Interval: 0,
		}
		rec.SetLocation(loc)

		result := NextForRecurrence(from, timeOfDay, rec)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		// Should advance by 1 day
		expected := from.Add(24 * time.Hour)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, *result)
		}
	})

	t.Run("Unknown recurrence type defaults to daily", func(t *testing.T) {
		rec := &entities.Recurrence{
			Type: entities.RecurrenceType(999), // Unknown type
		}
		rec.SetLocation(loc)

		result := NextForRecurrence(from, timeOfDay, rec)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		// Should behave like daily
		expected := NextDailyTrigger(from, timeOfDay, loc)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, *result)
		}
	})
}
