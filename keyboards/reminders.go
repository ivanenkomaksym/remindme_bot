package keyboards

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

const (
	CallbackRemindersList        = "rem_list"
	CallbackReminderDeletePrefix = "rem_del:"
)

func IsRemindersCallback(callbackData string) bool {
	return callbackData == CallbackRemindersList || strings.HasPrefix(callbackData, CallbackReminderDeletePrefix)
}

func GetRemindersListMarkup(reminders []entities.Reminder, lang string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	if len(reminders) == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu),
		))
		menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
		return &menu
	}

	for _, r := range reminders {
		// Add indicator for active/inactive status
		status := "❌" // inactive by default
		if r.IsActive {
			status = "✅" // active
		}

		label := fmt.Sprintf("%s %s", status, RecurrenceTypeLabel(lang, r.Recurrence.Type))
		// Add extra detail for monthly recurrences
		if r.Recurrence.IsMonthly() {
			// format days like: 1, 5, 12
			if len(r.Recurrence.DayOfMonth) > 0 {
				var daysStr strings.Builder
				for i, d := range r.Recurrence.DayOfMonth {
					if i > 0 {
						daysStr.WriteString(", ")
					}
					daysStr.WriteString(fmt.Sprintf("%d", d))
				}
				label = fmt.Sprintf("%s • %s", label, daysStr.String())
			}
		}
		// Always append time of day (or full date for once below)
		label = fmt.Sprintf("%s • %s", label, r.Recurrence.GetTimeOfDay())
		btn := tgbotapi.NewInlineKeyboardButtonData(
			s.BtnDelete,
			fmt.Sprintf("%s%d", CallbackReminderDeletePrefix, r.ID),
		)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, "noop"),
				btn,
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func FormatRemindersListText(reminders []entities.Reminder, lang string) string {
	s := T(lang)
	if len(reminders) == 0 {
		return s.NoReminders
	}
	var b strings.Builder
	b.WriteString(s.YourReminders)
	for _, r := range reminders {
		// Add indicator for active/inactive status
		status := "❌" // inactive by default
		if r.IsActive {
			status = "✅" // active
		}

		recurrenceType := RecurrenceTypeLabel(lang, r.Recurrence.Type)

		reminderTime := r.Recurrence.GetTimeOfDay()
		if r.Recurrence.Type == entities.Once {
			reminderTime = r.Recurrence.StartDate.Format("2006-01-02T15:04:05")
		}

		b.WriteString(fmt.Sprintf("%s • %s %s %s — %s (ID %d)\n", status, recurrenceType, s.At, reminderTime, r.Message, r.ID))
	}
	return b.String()
}

func ParseDeleteReminderID(callbackData string) (int64, bool) {
	if !strings.HasPrefix(callbackData, CallbackReminderDeletePrefix) {
		return 0, false
	}
	idStr := strings.TrimPrefix(callbackData, CallbackReminderDeletePrefix)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
