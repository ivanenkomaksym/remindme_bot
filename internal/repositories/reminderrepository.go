package repositories

import "github.com/ivanenkomaksym/remindme_bot/internal/models"

type ReminderRepository interface {
	CreateDailyReminder(time string, user models.User) *models.Reminder
	GetReminders() []models.Reminder
}
