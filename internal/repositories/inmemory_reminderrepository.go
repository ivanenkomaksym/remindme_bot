package repositories

import (
	"sync"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/scheduler"
)

type InMemoryReminderRepository struct {
	mu        sync.Mutex
	nextID    int64
	reminders []models.Reminder
}

func NewInMemoryReminderRepository() *InMemoryReminderRepository {
	return &InMemoryReminderRepository{
		nextID:    1,
		reminders: make([]models.Reminder, 0),
	}
}

func (r *InMemoryReminderRepository) CreateDailyReminder(timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.DailyAt(timeStr)
	next := scheduler.NextDailyTrigger(now, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
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

func (r *InMemoryReminderRepository) CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.CustomWeekly(daysOfWeek, timeStr)
	next := scheduler.NextWeeklyTrigger(now, daysOfWeek, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
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

func (r *InMemoryReminderRepository) CreateMonthlyReminder(daysOfMonth []int, timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.MonthlyOnDay(daysOfMonth, timeStr)
	next := scheduler.NextMonthlyTrigger(now, daysOfMonth, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
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

func (r *InMemoryReminderRepository) GetReminders() []models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]models.Reminder, len(r.reminders))
	copy(out, r.reminders)
	return out
}

func (r *InMemoryReminderRepository) GetRemindersByUser(userID int64) []models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]models.Reminder, 0)
	for _, rem := range r.reminders {
		if rem.User.Id == userID {
			result = append(result, rem)
		}
	}
	return result
}

func (r *InMemoryReminderRepository) DeleteReminder(reminderID int64, userID int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rem := range r.reminders {
		if rem.ID == reminderID && rem.User.Id == userID {
			// delete without preserving order
			r.reminders[i] = r.reminders[len(r.reminders)-1]
			r.reminders = r.reminders[:len(r.reminders)-1]
			return true
		}
	}
	return false
}

func (r *InMemoryReminderRepository) UpdateReminder(reminder *models.Reminder) bool {
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
