package notifier

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"
	"github.com/ivanenkomaksym/remindme_bot/internal/scheduler"
)

// BotSender is a minimal interface of the bot used for sending messages.
type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// startReminderNotifier runs a loop that checks for due reminders and notifies users
func StartReminderNotifier(reminderRepo repositories.ReminderRepository, bot *tgbotapi.BotAPI) {
	for {
		ProcessDueReminders(time.Now(), reminderRepo, bot)
		time.Sleep(time.Minute) // Check every minute
	}
}

// ProcessDueReminders performs a single pass over repository reminders, sending due ones
// and updating their next trigger. Extracted for testability.
func ProcessDueReminders(now time.Time, reminderRepo repositories.ReminderRepository, sender BotSender) {
	reminders := reminderRepo.GetReminders()
	for i := range reminders {
		rem := &reminders[i]
		if !rem.IsActive {
			continue
		}
		if rem.NextTrigger.After(now) {
			continue
		}

		// Send notification
		msg := tgbotapi.NewMessage(rem.User.Id, rem.Message)
		if _, err := sender.Send(msg); err != nil {
			log.Printf("Failed to send reminder to user %d: %v", rem.User.Id, err)
		}

		// Update NextTrigger for recurring reminders
		if rem.Recurrence != nil {
			next := scheduler.NextForRecurrence(rem.NextTrigger, rem.Recurrence)
			rem.NextTrigger = next
		} else {
			rem.IsActive = false // deactivate one-time reminders
		}
		reminderRepo.UpdateReminder(rem)
	}
}
