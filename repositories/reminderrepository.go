package repositories

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

type ReminderRepository interface {
	CreateDailyReminder(time string, user entities.User, message string) *entities.Reminder
	CreateWeeklyReminder(days_of_week []time.Weekday, time string, user entities.User, message string) *entities.Reminder
	CreateMonthlyReminder(days_of_month []int, time string, user entities.User, message string) *entities.Reminder
	GetReminders() []entities.Reminder
	GetRemindersByUser(userID int64) []entities.Reminder
	DeleteReminder(reminderID int64, userID int64) bool
	UpdateReminder(reminder *entities.Reminder) bool
}
