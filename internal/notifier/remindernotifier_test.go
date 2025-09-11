package notifier

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"
)

type fakeSender struct{ sent int }

func (f *fakeSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.sent++
	return tgbotapi.Message{}, nil
}

func TestProcessDueReminders_DailyAdvancesNextTrigger(t *testing.T) {
	repo := repositories.NewInMemoryReminderRepository()
	user := models.User{Id: 123}
	// Create a daily reminder at 00:00, set its NextTrigger to a known past time
	rem := repo.CreateDailyReminder("00:00", user, "ping")
	// Force NextTrigger to be at a fixed point
	past := time.Now().Add(-2 * time.Hour).Truncate(time.Minute)
	rem.NextTrigger = past
	repo.UpdateReminder(rem)

	sender := &fakeSender{}

	now := past.Add(1 * time.Minute)
	ProcessDueReminders(now, repo, sender)

	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}

	// After processing, NextTrigger should advance by 24h for daily recurrence
	updated := repo.GetReminders()[0]
	want := past.Add(24 * time.Hour)
	if !updated.NextTrigger.Equal(want) {
		t.Fatalf("expected NextTrigger %v, got %v", want, updated.NextTrigger)
	}
	if !updated.IsActive {
		t.Fatalf("daily reminder should remain active")
	}
}

func TestProcessDueReminders_OneTimeDeactivates(t *testing.T) {
	repo := repositories.NewInMemoryReminderRepository()
	// Manually insert one-time reminder (Recurrence=nil)
	now := time.Now()
	rem := models.Reminder{
		ID:          999,
		User:        models.User{Id: 42},
		Message:     "one-time",
		CreatedAt:   now.Add(-time.Hour),
		NextTrigger: now.Add(-time.Minute),
		Recurrence:  nil,
		IsActive:    true,
	}
	// Inject into repo via UpdateReminder path after appending
	// Use repository's internal behavior by creating a daily and replacing it
	junk := repo.CreateDailyReminder("00:00", rem.User, "junk")
	rem.ID = junk.ID
	repo.UpdateReminder(&rem)

	sender := &fakeSender{}

	ProcessDueReminders(now, repo, sender)

	if sender.sent != 1 {
		t.Fatalf("expected 1 message sent, got %d", sender.sent)
	}

	updated := repo.GetReminders()[0]
	if updated.IsActive {
		t.Fatalf("one-time reminder should be deactivated after sending")
	}
}
