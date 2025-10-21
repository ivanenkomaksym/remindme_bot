package notifier

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

type fakeSender struct{ sent int }

func (f *fakeSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.sent++
	return tgbotapi.Message{}, nil
}

func TestCalculateNextAlignedTime(t *testing.T) {
	tests := []struct {
		name     string
		now      string
		interval time.Duration
		expected string
	}{
		{
			name:     "9:28 with 15min interval should align to 9:30",
			now:      "09:28:00",
			interval: 15 * time.Minute,
			expected: "09:30:00",
		},
		{
			name:     "9:46 with 15min interval should align to 10:00",
			now:      "09:46:00",
			interval: 15 * time.Minute,
			expected: "10:00:00",
		},
		{
			name:     "9:46 with 5min interval should align to 9:50",
			now:      "09:46:00",
			interval: 5 * time.Minute,
			expected: "09:50:00",
		},
		{
			name:     "9:30 exact with 15min interval should align to 9:45",
			now:      "09:30:00",
			interval: 15 * time.Minute,
			expected: "09:45:00",
		},
		{
			name:     "23:55 with 15min interval should align to 00:00 next day",
			now:      "23:55:00",
			interval: 15 * time.Minute,
			expected: "00:00:00",
		},
		{
			name:     "10:17 with 10min interval should align to 10:20",
			now:      "10:17:00",
			interval: 10 * time.Minute,
			expected: "10:20:00",
		},
		{
			name:     "14:58 with 30min interval should align to 15:00",
			now:      "14:58:00",
			interval: 30 * time.Minute,
			expected: "15:00:00",
		},
		{
			name:     "23:45 with 30min interval should align to 00:00 next day",
			now:      "23:45:00",
			interval: 30 * time.Minute,
			expected: "00:00:00",
		},
	}

	baseDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the time strings relative to baseDate
			nowTime, err := time.Parse("15:04:05", tt.now)
			if err != nil {
				t.Fatalf("Failed to parse now time: %v", err)
			}
			now := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
				nowTime.Hour(), nowTime.Minute(), nowTime.Second(), 0, time.UTC)

			expectedTime, err := time.Parse("15:04:05", tt.expected)
			if err != nil {
				t.Fatalf("Failed to parse expected time: %v", err)
			}

			var expected time.Time
			if expectedTime.Hour() == 0 && expectedTime.Minute() == 0 &&
				(now.Hour() > 20 || (now.Hour() == 23 && now.Minute() >= 30)) {
				// Next day case
				nextDay := baseDate.Add(24 * time.Hour)
				expected = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
					expectedTime.Hour(), expectedTime.Minute(), expectedTime.Second(), 0, time.UTC)
			} else {
				// Same day case
				expected = time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
					expectedTime.Hour(), expectedTime.Minute(), expectedTime.Second(), 0, time.UTC)
			}

			result := calculateNextAlignedTime(now, tt.interval)

			if !result.Equal(expected) {
				t.Errorf("calculateNextAlignedTime(%s, %v) = %s, want %s",
					now.Format("15:04:05"), tt.interval, result.Format("15:04:05"), expected.Format("15:04:05"))
			}
		})
	}
}

func TestCalculateNextAlignedTime_EdgeCases(t *testing.T) {
	baseDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	t.Run("End of day overflow", func(t *testing.T) {
		// 23:50 with 15min interval should go to next day 00:00
		now := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 23, 50, 0, 0, time.UTC)
		expected := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day()+1, 0, 0, 0, 0, time.UTC)

		result := calculateNextAlignedTime(now, 15*time.Minute)

		if !result.Equal(expected) {
			t.Errorf("Expected next day 00:00, got %s", result.Format("2006-01-02 15:04:05"))
		}
	})

	t.Run("1 minute interval", func(t *testing.T) {
		now := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 12, 34, 30, 0, time.UTC)
		expected := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 12, 35, 0, 0, time.UTC)

		result := calculateNextAlignedTime(now, 1*time.Minute)

		if !result.Equal(expected) {
			t.Errorf("Expected 12:35:00, got %s", result.Format("15:04:05"))
		}
	})

	t.Run("60 minute interval", func(t *testing.T) {
		now := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 14, 30, 0, 0, time.UTC)
		expected := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 15, 0, 0, 0, time.UTC)

		result := calculateNextAlignedTime(now, 60*time.Minute)

		if !result.Equal(expected) {
			t.Errorf("Expected 15:00:00, got %s", result.Format("15:04:05"))
		}
	})
}

func TestCalculateNextAlignedTime_PreservesTimezone(t *testing.T) {
	// Test that the function preserves the input timezone
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("EST timezone not available")
	}

	now := time.Date(2025, 1, 15, 14, 28, 0, 0, est)
	result := calculateNextAlignedTime(now, 15*time.Minute)
	expected := time.Date(2025, 1, 15, 14, 30, 0, 0, est)

	if !result.Equal(expected) {
		t.Errorf("Expected %s, got %s", expected.Format("15:04:05 MST"), result.Format("15:04:05 MST"))
	}

	if result.Location() != est {
		t.Errorf("Expected timezone %s, got %s", est, result.Location())
	}
}

func TestProcessDueReminders_DailyAdvancesNextTrigger(t *testing.T) {
	repo := inmemory.NewInMemoryReminderRepository()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	user := entities.User{ID: 123, Location: loc}
	// Force NextTrigger to be at a fixed point in user's timezone
	past := time.Now().Add(-2 * time.Hour).Truncate(time.Minute).In(loc)
	// Create a daily reminder at the same clock time as 'past'
	tod := past
	rem, _ := repo.CreateDailyReminder(tod, &user, "ping")
	rem.NextTrigger = &past
	repo.UpdateReminder(rem)

	sender := &fakeSender{}

	now := past.Add(1 * time.Minute)
	ProcessDueReminders(now, repo, sender)

	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}

	// After processing, NextTrigger should advance according to scheduler.NextForRecurrence
	reminders, _ := repo.GetReminders()
	updated := reminders[0]
	expected := scheduler.NextForRecurrence(past, past, rem.Recurrence)
	if expected == nil {
		t.Fatalf("expected non-nil expected NextTrigger")
	}
	if updated.NextTrigger == nil {
		t.Fatalf("expected NextTrigger to be set, got nil")
	}
	// Compare instants with tolerance to account for time zone conversions and rounding
	abs := func(d time.Duration) time.Duration {
		if d < 0 {
			return -d
		}
		return d
	}
	if abs(updated.NextTrigger.Sub(*expected)) > time.Minute {
		t.Fatalf("expected NextTrigger approx %v, got %v", *expected, updated.NextTrigger)
	}
	if !updated.IsActive {
		t.Fatalf("daily reminder should remain active")
	}
}

func TestProcessDueReminders_OneTimeDeactivates(t *testing.T) {
	repo := inmemory.NewInMemoryReminderRepository()
	// Manually insert one-time reminder (Recurrence=nil)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	user := entities.User{ID: 123, Location: loc}
	now := time.Now()
	next := now.Add(-time.Minute)
	rem := entities.Reminder{
		ID:          999,
		UserID:      user.ID,
		Message:     "one-time",
		CreatedAt:   now.Add(-time.Hour),
		NextTrigger: &next,
		Recurrence:  nil,
		IsActive:    true,
	}
	// Inject into repo via UpdateReminder path after appending
	// Use repository's internal behavior by creating a daily and replacing it
	todJunk, _ := time.ParseInLocation("15:04", "00:00", loc)
	junk, _ := repo.CreateDailyReminder(todJunk, &user, "junk")
	rem.ID = junk.ID
	repo.UpdateReminder(&rem)

	sender := &fakeSender{}

	ProcessDueReminders(now, repo, sender)

	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}

	reminders, _ := repo.GetReminders()
	updated := reminders[0]
	if updated.IsActive {
		t.Fatalf("one-time reminder should be deactivated after sending")
	}
}

func TestProcessDueRemindersWithConfig_MonitoringDisabled(t *testing.T) {
	repo := inmemory.NewInMemoryReminderRepository()
	loc, _ := time.LoadLocation("UTC")
	user := entities.User{ID: 123, Location: loc}
	
	// Create a reminder that's due
	past := time.Now().Add(-2 * time.Hour).Truncate(time.Minute).In(loc)
	tod := past
	rem, _ := repo.CreateDailyReminder(tod, &user, "test message")
	rem.NextTrigger = &past
	repo.UpdateReminder(rem)

	sender := &fakeSender{}
	
	// Config with monitoring disabled
	botConfig := config.BotConfig{
		MonitorPendingUpdates:   false,
		PendingUpdatesThreshold: 100,
		AutoClearPendingUpdates: true,
	}

	now := past.Add(1 * time.Minute)
	ProcessDueRemindersWithConfig(now, repo, sender, botConfig)

	// Should still send reminder
	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}
	
	// Verify reminder was processed correctly
	reminders, _ := repo.GetReminders()
	updated := reminders[0]
	if !updated.IsActive {
		t.Fatalf("daily reminder should remain active")
	}
	if updated.NextTrigger == nil {
		t.Fatalf("NextTrigger should be updated")
	}
}

func TestProcessDueRemindersWithConfig_MonitoringEnabled(t *testing.T) {
	repo := inmemory.NewInMemoryReminderRepository()
	loc, _ := time.LoadLocation("UTC")
	user := entities.User{ID: 123, Location: loc}
	
	// Create a reminder that's due
	past := time.Now().Add(-2 * time.Hour).Truncate(time.Minute).In(loc)
	tod := past
	rem, _ := repo.CreateDailyReminder(tod, &user, "test message")
	rem.NextTrigger = &past
	repo.UpdateReminder(rem)

	sender := &fakeSender{}
	
	// Config with monitoring enabled
	botConfig := config.BotConfig{
		MonitorPendingUpdates:   true,
		PendingUpdatesThreshold: 50,
		AutoClearPendingUpdates: false,
	}

	now := past.Add(1 * time.Minute)
	ProcessDueRemindersWithConfig(now, repo, sender, botConfig)

	// Should still send reminder (monitoring doesn't affect core functionality)
	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}
	
	// Verify reminder was processed correctly
	reminders, _ := repo.GetReminders()
	updated := reminders[0]
	if !updated.IsActive {
		t.Fatalf("daily reminder should remain active")
	}
	if updated.NextTrigger == nil {
		t.Fatalf("NextTrigger should be updated")
	}
}
