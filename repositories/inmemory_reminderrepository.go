package repositories

import (
	"sync"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

type InMemoryReminderRepository struct {
	mu        sync.Mutex
	nextID    int64
	reminders []entities.Reminder
}

func NewInMemoryReminderRepository() *InMemoryReminderRepository {
	return &InMemoryReminderRepository{
		nextID:    1,
		reminders: make([]entities.Reminder, 0),
	}
}

func (r *InMemoryReminderRepository) CreateDailyReminder(timeStr string, user entities.User, message string) *entities.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := entities.DailyAt(timeStr)
	next := scheduler.NextDailyTrigger(now, timeStr)

	reminder := entities.Reminder{
		ID:          r.nextID,
		UserID:      user.ID,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user entities.User, message string) *entities.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := entities.CustomWeekly(daysOfWeek, timeStr)
	next := scheduler.NextWeeklyTrigger(now, daysOfWeek, timeStr)

	reminder := entities.Reminder{
		ID:          r.nextID,
		UserID:      user.ID,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) CreateMonthlyReminder(daysOfMonth []int, timeStr string, user entities.User, message string) *entities.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := entities.MonthlyOnDay(daysOfMonth, timeStr)
	next := scheduler.NextMonthlyTrigger(now, daysOfMonth, timeStr)

	reminder := entities.Reminder{
		ID:          r.nextID,
		UserID:      user.ID,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) GetReminders() []entities.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]entities.Reminder, len(r.reminders))
	copy(out, r.reminders)
	return out
}

func (r *InMemoryReminderRepository) GetRemindersByUser(userID int64) []entities.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]entities.Reminder, 0)
	for _, rem := range r.reminders {
		if rem.UserID == userID {
			result = append(result, rem)
		}
	}
	return result
}

func (r *InMemoryReminderRepository) DeleteReminder(reminderID int64, userID int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rem := range r.reminders {
		if rem.ID == reminderID && rem.UserID == userID {
			// delete without preserving order
			r.reminders[i] = r.reminders[len(r.reminders)-1]
			r.reminders = r.reminders[:len(r.reminders)-1]
			return true
		}
	}
	return false
}

func (r *InMemoryReminderRepository) UpdateReminder(reminder *entities.Reminder) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.reminders {
		if r.reminders[i].ID == reminder.ID {
			r.reminders[i] = *reminder
			return true
		}
	}
	return false
}
