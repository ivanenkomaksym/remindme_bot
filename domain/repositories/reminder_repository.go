package repositories

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// ReminderRepository defines the interface for reminder data operations
type ReminderRepository interface {
	// Reminder creation
	CreateDailyReminder(timeStr string, user *entities.User, message string) (*entities.Reminder, error)
	CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user *entities.User, message string) (*entities.Reminder, error)
	CreateMonthlyReminder(daysOfMonth []int, timeStr string, user *entities.User, message string) (*entities.Reminder, error)

	// Reminder retrieval
	GetReminders() ([]entities.Reminder, error)
	GetRemindersByUser(userID int64) ([]entities.Reminder, error)
	GetReminder(reminderID int64) (*entities.Reminder, error)

	// Reminder management
	UpdateReminder(reminder *entities.Reminder) error
	DeleteReminder(reminderID int64, userID int64) error
	DeactivateReminder(reminderID int64, userID int64) error

	// Reminder scheduling
	GetActiveReminders() ([]entities.Reminder, error)
	UpdateNextTrigger(reminderID int64, nextTrigger time.Time) error
}
