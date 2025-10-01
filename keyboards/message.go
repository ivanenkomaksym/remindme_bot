package keyboards

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// The callback data prefixes help to parse the user's message selection.
const (
	CallbackPrefixMessage = "msg_"
	CallbackMessageCustom = "msg_custom"
)

func IsMessageSelectionCallback(callbackData string) bool {
	return strings.HasPrefix(callbackData, CallbackPrefixMessage)
}

func GetMessageSelectionMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	// Add default message options
	for i, msg := range s.DefaultMessages {
		callbackData := CallbackPrefixMessage + string(rune(i))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(msg, callbackData)))
	}

	// Add custom message and confirm options
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ "+s.MsgEnterCustomMessage, CallbackMessageCustom),
	))

	// Add back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackTimeStart),
	))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func HandleMessageSelection(callbackData string,
	user *entities.User,
	userSelection *entities.UserSelection) (*SelectionResult, bool) {

	if callbackData == CallbackMessageCustom {
		userSelection.CustomText = true
		s := T(user.Language)
		return &SelectionResult{Text: s.MsgEnterCustomMessage, Markup: nil}, false
	}

	s := T(user.Language)
	if strings.HasPrefix(callbackData, CallbackPrefixMessage) {
		// Extract message index
		msgIndex := int(callbackData[len(CallbackPrefixMessage)])
		if msgIndex >= 0 && msgIndex < len(s.DefaultMessages) {
			userSelection.ReminderMessage = s.DefaultMessages[msgIndex]
		}
		return nil, true
	}

	return &SelectionResult{Text: s.MsgSelectMessage, Markup: GetMessageSelectionMarkup(user.Language)}, false
}

func HandleCustomText(text string,
	msg *tgbotapi.MessageConfig,
	user *entities.User,
	userSelection *entities.UserSelection) (*SelectionResult, bool) {
	userSelection.ReminderMessage = text
	return nil, true
}

func FormatReminderConfirmation(user *entities.User, userSelection *entities.UserSelection) *SelectionResult {
	s := T(user.Language)

	confirmation := "✅ " + s.ReminderSet + "!\n\n"
	confirmation += "📅 " + s.Frequency + ": " + RecurrenceTypeLabel(user.Language, userSelection.RecurrenceType) + "\n"

	if userSelection.RecurrenceType == entities.Weekly {
		confirmation += "📆 " + s.Days + ": "
		days := []string{}
		for _, weekday := range userSelection.WeekOptions {
			days = append(days, s.WeekdayNames[weekday])
		}
		if len(days) > 0 {
			confirmation += strings.Join(days, ", ")
		} else {
			confirmation += s.NoneSelected
		}
		confirmation += "\n"
	}

	if userSelection.RecurrenceType == entities.Monthly {
		confirmation += "📆 " + s.Days + ": "
		var days []string
		for i, selected := range userSelection.MonthOptions {
			if selected {
				days = append(days, fmt.Sprintf("%d", i+1))
			}
		}
		if len(days) > 0 {
			confirmation += strings.Join(days, ", ")
		} else {
			confirmation += s.NoneSelected
		}
		confirmation += "\n"
	}

	if userSelection.RecurrenceType == entities.Once {
		confirmation += "📅 " + s.Date + ": " + userSelection.SelectedDate.Format("2006-01-02") + "\n"
	}

	if userSelection.RecurrenceType == entities.Interval {
		if userSelection.IntervalDays > 0 {
			confirmation += "📆 " + fmt.Sprintf(s.MsgEveryNDays, userSelection.IntervalDays) + "\n"
		} else {
			confirmation += "📆 " + s.MsgIntervalPrompt + "\n"
		}
	}

	if userSelection.RecurrenceType == entities.SpacedBasedRepetition {
		var days string
		// TODO: make this more flexible in the future and reuse recurrence parsing logic
		for i, c := range []int{1, 2, 3, 5, 7, 7, 7} {
			if i > 0 {
				days += ", "
			}
			days += fmt.Sprintf("%d", c)
		}
		confirmation += "📆 " + fmt.Sprintf(s.MsgEveryNDaysSpaced, days) + "\n"
	}

	confirmation += "⏰ " + s.Time + ": " + userSelection.SelectedTime + "\n"
	confirmation += "💬 " + s.Message + ": " + userSelection.ReminderMessage + "\n\n"
	confirmation += s.ReminderScheduled

	myRemindersMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
		),
	)

	return &SelectionResult{Text: confirmation, Markup: &myRemindersMenu}
}
