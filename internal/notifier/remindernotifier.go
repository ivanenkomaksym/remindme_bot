package notifier

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"
)

// startReminderNotifier runs a loop that checks for due reminders and notifies users
func StartReminderNotifier(reminderRepo repositories.ReminderRepository, bot *tgbotapi.BotAPI) {
	for {
		now := time.Now()
		reminders := reminderRepo.GetReminders()
		for i := range reminders {
			rem := &reminders[i]
			if !rem.IsActive {
				continue
			}
			if rem.NextTrigger.Before(now) || rem.NextTrigger.Equal(now) {
				// Send notification
				msg := tgbotapi.NewMessage(rem.User.Id, rem.Message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Failed to send reminder to user %d: %v", rem.User.Id, err)
				}

				// Update NextTrigger for recurring reminders
				if rem.Recurrence != nil {
					next := getNextRecurrence(rem.NextTrigger, rem.Recurrence)
					rem.NextTrigger = next
				} else {
					rem.IsActive = false // deactivate one-time reminders
				}
				// TODO: persist the updated reminder (if using persistent storage)
			}
		}
		time.Sleep(time.Second * 10) // Check every 10 seconds
	}
}

// getNextRecurrence calculates the next trigger time for a recurring reminder
func getNextRecurrence(last time.Time, rec *models.Recurrence) time.Time {
	switch rec.Type {
	case models.Daily:
		return last.Add(24 * time.Hour)
	// TODO: Implement Weekly, Monthly, Interval and Custom types
	default:
		return last.Add(24 * time.Hour)
	}
}
