package notifier

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

type fakeSender struct{ sent int }

func (f *fakeSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.sent++
	return tgbotapi.Message{}, nil
}

func TestProcessDueReminders_DailyAdvancesNextTrigger(t *testing.T) {
	repo := inmemory.NewInMemoryReminderRepository()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	user := entities.User{ID: 123, Location: loc}
	// Force NextTrigger to be at a fixed point in user's timezone
	past := time.Now().Add(-2 * time.Hour).Truncate(time.Minute).In(loc)
	// Create a daily reminder at the same clock time as 'past'
	timeStr := past.Format("15:04")
	rem, _ := repo.CreateDailyReminder(timeStr, &user, "ping")
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
	expected := scheduler.NextForRecurrence(past, rem.Recurrence)
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
	junk, _ := repo.CreateDailyReminder("00:00", &user, "junk")
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
