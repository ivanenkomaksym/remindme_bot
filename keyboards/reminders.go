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
		label := fmt.Sprintf("%s • %s", RecurrenceTypeLabel(lang, r.Recurrence.Type), r.NextTrigger.Format("15:04"))
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
		b.WriteString(fmt.Sprintf("• %s %s %s — %s (ID %d)\n", RecurrenceTypeLabel(lang, r.Recurrence.Type), s.At, r.NextTrigger.Format("15:04"), r.Message, r.ID))
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
