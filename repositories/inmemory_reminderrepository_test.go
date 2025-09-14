package repositories

import (
	"slices"
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func TestCreateDailyReminder_Happy(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 1, UserName: "tester"}
	rem := repo.CreateDailyReminder("23:15", user, "daily msg")

	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	if rem.ID != 1 {
		t.Errorf("expected ID 1, got %d", rem.ID)
	}
	if !rem.IsActive {
		t.Errorf("expected IsActive true")
	}
	if !rem.Recurrence.IsDaily() {
		t.Errorf("expected daily recurrence")
	}
	if rem.Recurrence.TimeOfDay != "23:15" {
		t.Errorf("expected TimeOfDay 23:15, got %s", rem.Recurrence.TimeOfDay)
	}
	if !rem.NextTrigger.After(rem.CreatedAt) {
		t.Errorf("expected NextTrigger after CreatedAt")
	}
	if rem.NextTrigger.Hour() != 23 || rem.NextTrigger.Minute() != 15 {
		t.Errorf("expected next at 23:15, got %02d:%02d", rem.NextTrigger.Hour(), rem.NextTrigger.Minute())
	}
}

func TestCreateDailyReminder_InvalidTime(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 2, UserName: "tester2"}
	rem := repo.CreateDailyReminder("bad", user, "daily bad")

	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	if !rem.NextTrigger.Equal(rem.CreatedAt) {
		t.Errorf("expected NextTrigger == CreatedAt on invalid time; got %v vs %v", rem.NextTrigger, rem.CreatedAt)
	}
	if !rem.Recurrence.IsDaily() {
		t.Errorf("expected daily recurrence even for invalid time string")
	}
}

func TestCreateWeeklyReminder_Happy(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 3}
	days := []time.Weekday{time.Wednesday, time.Friday}
	rem := repo.CreateWeeklyReminder(days, "00:01", user, "weekly msg")

	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	if rem.Recurrence == nil || !rem.Recurrence.IsWeekly() {
		t.Fatalf("expected weekly recurrence")
	}
	wd := rem.NextTrigger.Weekday()
	if !slices.Contains(days, wd) {
		t.Errorf("expected weekday in %v, got %v", days, wd)
	}
	if rem.NextTrigger.Sub(rem.CreatedAt) < 0 || rem.NextTrigger.Sub(rem.CreatedAt) > 7*24*time.Hour {
		t.Errorf("expected next trigger within 7 days; delta=%v", rem.NextTrigger.Sub(rem.CreatedAt))
	}
	if rem.NextTrigger.Hour() != 0 || rem.NextTrigger.Minute() != 1 {
		t.Errorf("expected next at 00:01, got %02d:%02d", rem.NextTrigger.Hour(), rem.NextTrigger.Minute())
	}
}

func TestCreateWeeklyReminder_EmptyDaysFallbackDaily(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 4}
	rem := repo.CreateWeeklyReminder([]time.Weekday{}, "06:30", user, "weekly empty")

	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	// Even though recurrence is weekly, logic falls back to daily scheduler; we only check timing behavior
	delta := rem.NextTrigger.Sub(rem.CreatedAt)
	if delta < 0 || delta > 24*time.Hour {
		t.Errorf("expected next trigger within 24h; delta=%v", delta)
	}
	if rem.NextTrigger.Hour() != 6 || rem.NextTrigger.Minute() != 30 {
		t.Errorf("expected next at 06:30, got %02d:%02d", rem.NextTrigger.Hour(), rem.NextTrigger.Minute())
	}
}

func TestCreateMonthlyReminder_Happy(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 5}
	days := []int{5, 20}
	rem := repo.CreateMonthlyReminder(days, "07:45", user, "monthly msg")
	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	if rem.Recurrence == nil || !rem.Recurrence.IsMonthly() {
		t.Fatalf("expected monthly recurrence")
	}
	d := rem.NextTrigger.Day()
	if !slices.Contains(days, d) {
		t.Errorf("expected day in %d, got %d", days, d)
	}
	if rem.NextTrigger.Sub(rem.CreatedAt) < 0 || rem.NextTrigger.Sub(rem.CreatedAt) > 35*24*time.Hour {
		t.Errorf("expected next trigger within ~35 days; delta=%v", rem.NextTrigger.Sub(rem.CreatedAt))
	}
	if rem.NextTrigger.Hour() != 7 || rem.NextTrigger.Minute() != 45 {
		t.Errorf("expected next at 07:45, got %02d:%02d", rem.NextTrigger.Hour(), rem.NextTrigger.Minute())
	}
}

func TestCreateMonthlyReminder_InvalidDaysFallbackDaily(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 6}
	rem := repo.CreateMonthlyReminder([]int{0, 35}, "09:00", user, "monthly invalid")
	if rem == nil {
		t.Fatalf("expected reminder, got nil")
	}
	delta := rem.NextTrigger.Sub(rem.CreatedAt)
	if delta < 0 || delta > 24*time.Hour {
		t.Errorf("expected next trigger within 24h for invalid days; delta=%v", delta)
	}
	if rem.NextTrigger.Hour() != 9 || rem.NextTrigger.Minute() != 0 {
		t.Errorf("expected next at 09:00, got %02d:%02d", rem.NextTrigger.Hour(), rem.NextTrigger.Minute())
	}
}

func TestGetReminders_ReturnsCopy(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 7}
	_ = repo.CreateDailyReminder("10:00", user, "original")
	list1 := repo.GetReminders()
	if len(list1) != 1 {
		t.Fatalf("expected 1 reminder, got %d", len(list1))
	}
	// mutate returned slice
	list1[0].Message = "changed"
	list2 := repo.GetReminders()
	if list2[0].Message != "original" {
		t.Errorf("expected internal data unchanged, got %q", list2[0].Message)
	}
}

func TestUpdateReminder_Happy(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 8}
	rem := repo.CreateDailyReminder("12:00", user, "original message")
	rem.Message = "updated message"
	rem.IsActive = false
	ok := repo.UpdateReminder(rem)
	if !ok {
		t.Fatalf("expected update to succeed")
	}
	updated := repo.GetReminders()[0]
	if updated.Message != "updated message" {
		t.Errorf("expected message 'updated message', got %q", updated.Message)
	}
	if updated.IsActive != false {
		t.Errorf("expected IsActive false, got %v", updated.IsActive)
	}
}

func TestUpdateReminder_NotFound(t *testing.T) {
	repo := NewInMemoryReminderRepository()
	user := entities.User{ID: 9}
	rem := &entities.Reminder{ID: 999, UserID: user.ID, Message: "does not exist"}
	ok := repo.UpdateReminder(rem)
	if ok {
		t.Fatalf("expected update to fail for non-existent reminder")
	}
}
