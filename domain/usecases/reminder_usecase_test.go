package usecases

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
)

func newReminderUC() ReminderUseCase {
	remRepo := inmemory.NewInMemoryReminderRepository()
	userRepo := inmemory.NewInMemoryUserRepository()
	return NewReminderUseCase(remRepo, userRepo)
}

func TestCreateReminder_ValidOnce(t *testing.T) {
	uc := newReminderUC()
	// seed user
	userRepo := inmemory.NewInMemoryUserRepository()
	userRepo.CreateOrUpdateUser(1, "u", "f", "l", "en")

	// Build a new UC that shares the same user repo as above
	remRepo := inmemory.NewInMemoryReminderRepository()
	uc = NewReminderUseCase(remRepo, userRepo)

	sel := entities.NewUserSelection()
	sel.RecurrenceType = entities.Once
	sel.SelectedDate = time.Now().Add(24 * time.Hour) // Tomorrow
	sel.SelectedTime = "12:00"
	sel.ReminderMessage = "Check e-mail"

	rem, err := uc.CreateReminder(1, sel)
	if rem.NextTrigger == nil {
		t.Fatalf("expected Next trigger to be set for once reminder")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
}

func TestCreateReminder_ValidDaily(t *testing.T) {
	uc := newReminderUC()
	// seed user
	userRepo := inmemory.NewInMemoryUserRepository()
	userRepo.CreateOrUpdateUser(1, "u", "f", "l", "en")

	// Build a new UC that shares the same user repo as above
	remRepo := inmemory.NewInMemoryReminderRepository()
	uc = NewReminderUseCase(remRepo, userRepo)

	sel := entities.NewUserSelection()
	sel.RecurrenceType = entities.Daily
	sel.SelectedTime = "09:00"
	sel.ReminderMessage = "Take a break"

	rem, err := uc.CreateReminder(1, sel)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
}

func TestCreateReminder_ValidationErrors(t *testing.T) {
	uc := newReminderUC()

	// nil selection
	if _, err := uc.CreateReminder(1, nil); err == nil {
		t.Fatalf("expected error for nil selection")
	}

	// missing message
	sel := entities.NewUserSelection()
	sel.RecurrenceType = entities.Daily
	sel.SelectedTime = "10:00"
	if _, err := uc.CreateReminder(1, sel); err == nil {
		t.Fatalf("expected error for empty message")
	}

	// missing time
	sel = entities.NewUserSelection()
	sel.RecurrenceType = entities.Daily
	sel.ReminderMessage = "msg"
	if _, err := uc.CreateReminder(1, sel); err == nil {
		t.Fatalf("expected error for invalid time")
	}
}

func TestDeleteReminder_Validations(t *testing.T) {
	uc := newReminderUC()

	if err := uc.DeleteReminder(0, 1); err == nil {
		t.Fatalf("expected error for invalid reminder id")
	}
	if err := uc.DeleteReminder(1, 0); err == nil {
		t.Fatalf("expected error for invalid user id")
	}
}

func TestGetUserReminders_InvalidUser(t *testing.T) {
	uc := newReminderUC()
	if _, err := uc.GetUserReminders(0); err == nil {
		t.Fatalf("expected error for invalid user id")
	}
}

func TestCreateReminder_UserMustExist(t *testing.T) {
	uc := newReminderUC()
	sel := entities.NewUserSelection()
	sel.RecurrenceType = entities.Daily
	sel.SelectedTime = "09:00"
	sel.ReminderMessage = "Hello"
	_, err := uc.CreateReminder(999, sel)
	if err == nil || err != errors.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
