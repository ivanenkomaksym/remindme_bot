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
		label := formatLabel(r, lang, false)
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
		label := formatLabel(r, lang, true)
		b.WriteString(fmt.Sprintf("%s\n", label))
	}
	return b.String()
}

func formatLabel(reminder entities.Reminder, lang string, includeMessage bool) string {
	s := T(lang)

	// Add indicator for active/inactive status
	status := "❌" // inactive by default
	if reminder.IsActive {
		status = "✅" // active
	}

	var label string
	recurrenceType := RecurrenceTypeLabel(lang, reminder.Recurrence.Type)
	reminderTime := reminder.Recurrence.GetTimeOfDay()

	switch reminder.Recurrence.Type {
	case entities.Once:
		reminderTime = reminder.Recurrence.StartDate.In(reminder.Recurrence.Location).Format("2006-01-02T15:04:05")
	case entities.Daily:
		break
	case entities.Weekly:
		// format days like: 1, 5, 12
		var daysStr strings.Builder
		for i, d := range reminder.Recurrence.Weekdays {
			if i > 0 {
				daysStr.WriteString(", ")
			}
			daysStr.WriteString(fmt.Sprintf("%s", s.WeekdayNamesShort[d]))
		}
		reminderTime = fmt.Sprintf("%s • %s", daysStr.String(), reminderTime)
	case entities.Monthly:
		// format days like: 1, 5, 12
		var daysStr strings.Builder
		for i, d := range reminder.Recurrence.DayOfMonth {
			if i > 0 {
				daysStr.WriteString(", ")
			}
			daysStr.WriteString(fmt.Sprintf("%d", d))
		}
		reminderTime = fmt.Sprintf("%s • %s", daysStr.String(), reminderTime)
	case entities.Interval:
		reminderTime = fmt.Sprintf(s.MsgEveryNDays, reminder.Recurrence.Interval)
	}

	label = fmt.Sprintf("%s %s %s %s", status, recurrenceType, s.At, reminderTime)

	if includeMessage {
		label = fmt.Sprintf("%s — %s", label, reminder.Message)
	}

	return label
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
