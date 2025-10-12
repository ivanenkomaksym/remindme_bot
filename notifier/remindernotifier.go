package notifier

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

// BotSender is a minimal interface of the bot used for sending messages.
type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// startReminderNotifier runs a loop that checks for due reminders and notifies users
func StartReminderNotifier(reminderRepo repositories.ReminderRepository, appConfig config.AppConfig, bot *tgbotapi.BotAPI) {
	// Calculate the next aligned time to start
	now := time.Now()
	nextStart := calculateNextAlignedTime(now, appConfig.NotifierTimeout)

	// Wait until the aligned time
	initialDelay := nextStart.Sub(now)
	log.Printf("Starting reminder notifier in %v (next aligned time: %v)", initialDelay.Truncate(time.Second), nextStart.Format("15:04:05"))
	time.Sleep(initialDelay)

	// Now run the regular loop at aligned intervals
	for {
		ProcessDueReminders(time.Now(), reminderRepo, bot)
		time.Sleep(appConfig.NotifierTimeout)
	}
}

// calculateNextAlignedTime calculates the next time aligned to the interval
// For example, if interval is 15min and current time is 9:28, it returns 9:30
func calculateNextAlignedTime(now time.Time, interval time.Duration) time.Time {
	// Get minutes since midnight
	minutesSinceMidnight := now.Hour()*60 + now.Minute()
	intervalMinutes := int(interval.Minutes())

	// Calculate next aligned minute
	nextAlignedMinute := ((minutesSinceMidnight / intervalMinutes) + 1) * intervalMinutes

	// Handle day overflow
	if nextAlignedMinute >= 24*60 {
		// Next day
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
	}

	// Same day
	nextHour := nextAlignedMinute / 60
	nextMinute := nextAlignedMinute % 60

	return time.Date(now.Year(), now.Month(), now.Day(), nextHour, nextMinute, 0, 0, now.Location())
}

// ProcessDueReminders performs a single pass over repository reminders, sending due ones
// and updating their next trigger. Extracted for testability.
func ProcessDueReminders(now time.Time, reminderRepo repositories.ReminderRepository, sender BotSender) {
	reminders, _ := reminderRepo.GetReminders()
	for i := range reminders {
		rem := &reminders[i]
		if !rem.IsActive {
			continue
		}
		if rem.NextTrigger == nil || rem.NextTrigger.After(now) {
			continue
		}

		// Send notification
		msg := tgbotapi.NewMessage(rem.UserID, rem.Message)
		if _, err := sender.Send(msg); err != nil {
			log.Printf("Failed to send reminder to user %d: %v", rem.UserID, err)
		}

		// Update NextTrigger for recurring reminders
		if rem.Recurrence != nil && rem.Recurrence.Type != entities.Once {
			// Use StartDate for the time of day, not the previous NextTrigger
			timeOfDay := *rem.Recurrence.StartDate
			rem.NextTrigger = scheduler.NextForRecurrence(now, timeOfDay, rem.Recurrence)
		} else {
			rem.IsActive = false // deactivate one-time reminders
		}
		reminderRepo.UpdateReminder(rem)
	}
}
