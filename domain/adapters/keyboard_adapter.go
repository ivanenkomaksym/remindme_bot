package adapters

import (
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Conversion helpers removed: use entities.User and entities.UserSelection directly

// ConvertRemindersToModel is obsolete after migration to entities. Remove or refactor as needed.


// HandleRecurrenceTypeSelection is an adapter for keyboards.HandleRecurrenceTypeSelection
func HandleRecurrenceTypeSelection(callbackData string, user *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	userModel := user
	selectionModel := selection
	result, err := keyboards.HandleRecurrenceTypeSelection(callbackData, userModel, selectionModel)
	if err != nil {
		return nil, err
	}
	selection.RecurrenceType = selectionModel.RecurrenceType
	selection.IsWeekly = selectionModel.IsWeekly
	return result, nil
}

// HandleTimeSelection is an adapter for keyboards.HandleTimeSelection
func HandleTimeSelection(callbackData string, user *entities.User, selection *entities.UserSelection) *keyboards.SelectionResult {
	result := keyboards.HandleTimeSelection(callbackData, user, selection)
	selection.SelectedTime = selection.SelectedTime
	selection.CustomTime = selection.CustomTime
	return result
}

// HandleMessageSelection is an adapter for keyboards.HandleMessageSelection
func HandleMessageSelection(callbackData string, user *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, bool) {
	result, completed := keyboards.HandleMessageSelection(callbackData, user, selection)
	selection.ReminderMessage = selection.ReminderMessage
	selection.CustomText = selection.CustomText
	return result, completed
}

// HandleWeekSelection is an adapter for keyboards.HandleWeekSelection
func HandleWeekSelection(callbackData string, weekOptions *[7]bool, language string) *keyboards.SelectionResult {
	return keyboards.HandleWeekSelection(callbackData, weekOptions, language)
}

// FormatRemindersListText is an adapter for keyboards.FormatRemindersListText
func FormatRemindersListText(reminders []entities.Reminder, language string) string {
	// TODO: Refactor keyboards.FormatRemindersListText to accept entities.Reminder directly
	return "[Refactor required: FormatRemindersListText]"
}

// GetRemindersListMarkup is an adapter for keyboards.GetRemindersListMarkup
func GetRemindersListMarkup(reminders []entities.Reminder, language string) *tgbotapi.InlineKeyboardMarkup {
	// TODO: Refactor keyboards.GetRemindersListMarkup to accept entities.Reminder directly
	return nil
}

// HandleCustomTimeSelection is an adapter for keyboards.HadleCustomTimeSelection
func HandleCustomTimeSelection(text string, msg *tgbotapi.MessageConfig, user *entities.User, selection *entities.UserSelection) *keyboards.SelectionResult {
	selectionResult := keyboards.HadleCustomTimeSelection(text, msg, user, selection)
	// The selection is updated in place
	return selectionResult
}

// HandleCustomText is an adapter for keyboards.HadleCustomText
func HandleCustomText(text string, msg *tgbotapi.MessageConfig, user *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, bool) {
	selectionResult, completed := keyboards.HadleCustomText(text, msg, user, selection)
	// The selection is updated in place
	return selectionResult, completed
}
