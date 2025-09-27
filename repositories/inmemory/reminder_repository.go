package inmemory

import (
	"sync"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

type InMemoryReminderRepository struct {
	mu        sync.RWMutex
	nextID    int64
	reminders []entities.Reminder
}

func NewInMemoryReminderRepository() repositories.ReminderRepository {
	return &InMemoryReminderRepository{
		nextID:    1,
		reminders: make([]entities.Reminder, 0),
	}
}

// Reminder creation methods
func (r *InMemoryReminderRepository) CreateOnceReminder(date time.Time, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	recurrence := entities.OnceAt(date, timeStr, user.GetLocation())
	nextTrigger := *recurrence.StartDate
	reminder := entities.NewReminder(r.nextID, user.ID, message, recurrence, &nextTrigger)
	r.nextID++
	r.reminders = append(r.reminders, *reminder)

	return &r.reminders[len(r.reminders)-1], nil
}

func (r *InMemoryReminderRepository) CreateDailyReminder(timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	next := scheduler.NextDailyTrigger(now, timeStr, user.GetLocation())

	recurrence := entities.DailyAt(timeStr, user.GetLocation())
	reminder := entities.NewReminder(r.nextID, user.ID, message, recurrence, &next)
	r.nextID++
	r.reminders = append(r.reminders, *reminder)

	return &r.reminders[len(r.reminders)-1], nil
}

func (r *InMemoryReminderRepository) CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	next := scheduler.NextWeeklyTrigger(now, daysOfWeek, timeStr, user.GetLocation())
	// TODO: calculate next during creation of the reminder
	recurrence := entities.CustomWeekly(daysOfWeek, timeStr, user.GetLocation())
	reminder := entities.NewReminder(r.nextID, user.ID, message, recurrence, &next)
	r.nextID++
	r.reminders = append(r.reminders, *reminder)

	return &r.reminders[len(r.reminders)-1], nil
}

func (r *InMemoryReminderRepository) CreateMonthlyReminder(daysOfMonth []int, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	next := scheduler.NextMonthlyTrigger(now, daysOfMonth, timeStr, user.GetLocation())

	recurrence := entities.MonthlyOnDay(daysOfMonth, timeStr, user.GetLocation())
	reminder := entities.NewReminder(r.nextID, user.ID, message, recurrence, &next)
	r.nextID++
	r.reminders = append(r.reminders, *reminder)

	return &r.reminders[len(r.reminders)-1], nil
}

func (r *InMemoryReminderRepository) CreateIntervalReminder(intervalDays int, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	// Next trigger is N days from now at specified time
	base := scheduler.NextDailyTrigger(now, timeStr, user.GetLocation())
	next := base.Add(time.Duration(intervalDays-1) * 24 * time.Hour)

	recurrence := entities.IntervalEveryDays(intervalDays, timeStr, user.GetLocation())
	reminder := entities.NewReminder(r.nextID, user.ID, message, recurrence, &next)
	r.nextID++
	r.reminders = append(r.reminders, *reminder)

	return &r.reminders[len(r.reminders)-1], nil
}

// Reminder retrieval methods
func (r *InMemoryReminderRepository) GetReminders() ([]entities.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]entities.Reminder, len(r.reminders))
	copy(out, r.reminders)
	return out, nil
}

func (r *InMemoryReminderRepository) GetRemindersByUser(userID int64) ([]entities.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entities.Reminder, 0)
	for _, rem := range r.reminders {
		if rem.UserID == userID {
			result = append(result, rem)
		}
	}
	return result, nil
}

func (r *InMemoryReminderRepository) GetReminder(reminderID int64) (*entities.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, rem := range r.reminders {
		if rem.ID == reminderID {
			remCopy := rem
			return &remCopy, nil
		}
	}
	return nil, nil
}

// Reminder management methods
func (r *InMemoryReminderRepository) UpdateReminder(reminder *entities.Reminder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.reminders {
		if r.reminders[i].ID == reminder.ID {
			r.reminders[i] = *reminder
			return nil
		}
	}
	return nil // Reminder not found, nothing to update
}

func (r *InMemoryReminderRepository) DeleteReminder(reminderID int64, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, rem := range r.reminders {
		if rem.ID == reminderID && rem.UserID == userID {
			// Delete without preserving order
			r.reminders[i] = r.reminders[len(r.reminders)-1]
			r.reminders = r.reminders[:len(r.reminders)-1]
			return nil
		}
	}
	return nil // Reminder not found, nothing to delete
}

func (r *InMemoryReminderRepository) DeactivateReminder(reminderID int64, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.reminders {
		if r.reminders[i].ID == reminderID && r.reminders[i].UserID == userID {
			r.reminders[i].Deactivate()
			return nil
		}
	}
	return nil // Reminder not found, nothing to deactivate
}

// Reminder scheduling methods
func (r *InMemoryReminderRepository) GetActiveReminders() ([]entities.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entities.Reminder, 0)
	for _, rem := range r.reminders {
		if rem.IsActive {
			result = append(result, rem)
		}
	}
	return result, nil
}

func (r *InMemoryReminderRepository) UpdateNextTrigger(reminderID int64, nextTrigger time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.reminders {
		if r.reminders[i].ID == reminderID {
			r.reminders[i].UpdateNextTrigger(&nextTrigger)
			return nil
		}
	}
	return nil // Reminder not found, nothing to update
}
