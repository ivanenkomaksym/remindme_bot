package keyboards

import (
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
	confirmation += "📅 " + s.Frequency + ": " + userSelection.RecurrenceType.String() + "\n"

	if userSelection.RecurrenceType == entities.Weekly {
		confirmation += "📆 " + s.Days + ": "
		days := []string{}
		weekdayNames := s.WeekdayNames
		for i, selected := range userSelection.WeekOptions {
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
