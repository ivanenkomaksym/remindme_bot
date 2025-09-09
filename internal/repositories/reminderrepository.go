package repositories

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

type ReminderRepository interface {
	CreateDailyReminder(time string, user models.User, message string) *models.Reminder
	CreateWeeklyReminder(days_of_week []time.Weekday, time string, user models.User, message string) *models.Reminder
	CreateMonthlyReminder(days_of_month []int, time string, user models.User, message string) *models.Reminder
	GetReminders() []models.Reminder
	GetRemindersByUser(userID int64) []models.Reminder
	DeleteReminder(reminderID int64, userID int64) bool
}
