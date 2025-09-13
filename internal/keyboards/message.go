package keyboards

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
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
	msg *tgbotapi.EditMessageTextConfig,
	userState *types.UserSelectionState) (*tgbotapi.InlineKeyboardMarkup, bool) {

	if callbackData == CallbackMessageCustom {
		userState.CustomText = true
		s := T(userState.Language)
		msg.Text = s.MsgEnterCustomMessage
		return nil, false
	}

	s := T(userState.Language)
	if strings.HasPrefix(callbackData, CallbackPrefixMessage) {
		// Extract message index
		msgIndex := int(callbackData[len(CallbackPrefixMessage)])
		if msgIndex >= 0 && msgIndex < len(s.DefaultMessages) {
			userState.ReminderMessage = s.DefaultMessages[msgIndex]
		}
		return nil, true
	}

	return GetMessageSelectionMarkup(userState.Language), false
}

func HadleCustomText(text string,
	msg *tgbotapi.MessageConfig,
	userState *types.UserSelectionState) (*tgbotapi.InlineKeyboardMarkup, bool) {
	userState.ReminderMessage = text
	return nil, true
}

func FormatReminderConfirmation(userState *types.UserSelectionState) (string, *tgbotapi.InlineKeyboardMarkup) {
	s := T(userState.Language)

	confirmation := "✅ " + s.ReminderSet + "!\n\n"
	confirmation += "📅 " + s.Frequency + ": " + userState.RecurrenceType.String() + "\n"

	if userState.IsWeekly {
		confirmation += "📆 " + s.Days + ": "
		days := []string{}
		weekdayNames := s.WeekdayNames
		for i, selected := range userState.WeekOptions {
			if selected {
				days = append(days, weekdayNames[i])
			}
		}
		if len(days) > 0 {
			confirmation += strings.Join(days, ", ")
		} else {
			confirmation += s.NoneSelected
		}
		confirmation += "\n"
	}

	confirmation += "⏰ " + s.Time + ": " + userState.SelectedTime + "\n"
	confirmation += "💬 " + s.Message + ": " + userState.ReminderMessage + "\n\n"
	confirmation += s.ReminderScheduled

	myRemindersMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
		),
	)

	return confirmation, &myRemindersMenu
}
