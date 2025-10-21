package notifier

import (
	"fmt"
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
func StartReminderNotifier(reminderRepo repositories.ReminderRepository, appConfig config.AppConfig, botConfig config.BotConfig, bot *tgbotapi.BotAPI) {
	// Calculate the next aligned time to start
	now := time.Now()
	nextStart := calculateNextAlignedTime(now, appConfig.NotifierTimeout)

	// Wait until the aligned time
	initialDelay := nextStart.Sub(now)
	log.Printf("Starting reminder notifier in %v (next aligned time: %v)", initialDelay.Truncate(time.Second), nextStart.Format("15:04:05"))
	time.Sleep(initialDelay)

	// Now run the regular loop at aligned intervals
	for {
		ProcessDueRemindersWithConfig(time.Now(), reminderRepo, bot, botConfig)
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
	// Monitor Telegram bot pending updates if sender is the actual bot (with default config)
	if bot, ok := sender.(*tgbotapi.BotAPI); ok {
		defaultBotConfig := config.BotConfig{
			MonitorPendingUpdates:   true,
			PendingUpdatesThreshold: 100,
			AutoClearPendingUpdates: true,
		}
		monitorBotUpdatesWithConfig(bot, defaultBotConfig)
	}

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
		text := fmt.Sprintf("🔔 %s", rem.Message)
		msg := tgbotapi.NewMessage(rem.UserID, text)
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

// ProcessDueRemindersWithConfig performs a single pass over repository reminders with bot config
func ProcessDueRemindersWithConfig(now time.Time, reminderRepo repositories.ReminderRepository, sender BotSender, botConfig config.BotConfig) {
	// Monitor Telegram bot pending updates if sender is the actual bot and monitoring is enabled
	if bot, ok := sender.(*tgbotapi.BotAPI); ok && botConfig.MonitorPendingUpdates {
		monitorBotUpdatesWithConfig(bot, botConfig)
	}

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
		text := fmt.Sprintf("🔔 %s", rem.Message)
		msg := tgbotapi.NewMessage(rem.UserID, text)
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

// monitorBotUpdates checks for pending updates and logs them for debugging (with default config)
func monitorBotUpdates(bot *tgbotapi.BotAPI) {
	defaultBotConfig := config.BotConfig{
		MonitorPendingUpdates:   true,
		PendingUpdatesThreshold: 100,
		AutoClearPendingUpdates: true,
	}
	monitorBotUpdatesWithConfig(bot, defaultBotConfig)
}

// monitorBotUpdatesWithConfig checks for pending updates and logs them for debugging
func monitorBotUpdatesWithConfig(bot *tgbotapi.BotAPI, botConfig config.BotConfig) {
	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Printf("ERROR: Failed to get webhook info: %v", err)
		return
	}

	pendingUpdates := webhookInfo.PendingUpdateCount
	if pendingUpdates > 0 {
		log.Printf("WARNING: Bot has %d pending updates", pendingUpdates)

		// If there are too many pending updates, consider clearing them
		if pendingUpdates > botConfig.PendingUpdatesThreshold {
			log.Printf("CRITICAL: %d pending updates detected (threshold: %d), this may indicate service issues",
				pendingUpdates, botConfig.PendingUpdatesThreshold)

			// Optionally clear pending updates by setting webhook again
			if botConfig.AutoClearPendingUpdates {
				clearPendingUpdates(bot)
			} else {
				log.Printf("INFO: Auto-clear disabled, manual intervention required")
			}
		}
	} else {
		log.Printf("INFO: Bot status healthy - 0 pending updates")
	}

	// Log webhook status for debugging
	if webhookInfo.LastErrorDate != 0 {
		log.Printf("WARNING: Last webhook error at %s: %s",
			time.Unix(int64(webhookInfo.LastErrorDate), 0).Format("2006-01-02 15:04:05"),
			webhookInfo.LastErrorMessage)
	}
}

// clearPendingUpdates attempts to clear pending updates by resetting the webhook
func clearPendingUpdates(bot *tgbotapi.BotAPI) {
	log.Printf("INFO: Attempting to clear pending updates...")

	// Get current webhook info to preserve the URL
	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Printf("ERROR: Failed to get current webhook info: %v", err)
		return
	}

	// Set webhook with drop_pending_updates=true
	webhookConfig, err := tgbotapi.NewWebhook(webhookInfo.URL)
	if err != nil {
		log.Printf("ERROR: Failed to create webhook config: %v", err)
		return
	}
	webhookConfig.DropPendingUpdates = true

	_, err = bot.Request(webhookConfig)
	if err != nil {
		log.Printf("ERROR: Failed to clear pending updates: %v", err)
	} else {
		log.Printf("INFO: Successfully cleared pending updates")
	}
}
