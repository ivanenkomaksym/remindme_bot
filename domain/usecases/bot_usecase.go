package usecases

import (
	"log"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotUseCase defines the interface for bot business logic
type BotUseCase interface {
	HandleStartCommand(user *tgbotapi.User) (*keyboards.SelectionResult, error)
	HandleCallbackQuery(user *tgbotapi.User, callbackData string) (*keyboards.SelectionResult, error)
	HandleTextMessage(user *tgbotapi.User, text string) (*keyboards.SelectionResult, error)
	ProcessKeyboardSelection(callbackQuery *tgbotapi.CallbackQuery) (*keyboards.SelectionResult, error)
	ProcessUserInput(message *tgbotapi.Message) (*keyboards.SelectionResult, error)
}

type botUseCase struct {
	userUseCase     UserUseCase
	reminderUseCase ReminderUseCase
	bot             *tgbotapi.BotAPI
}

// NewBotUseCase creates a new bot use case
func NewBotUseCase(userUseCase UserUseCase, reminderUseCase ReminderUseCase, bot *tgbotapi.BotAPI) BotUseCase {
	return &botUseCase{
		userUseCase:     userUseCase,
		reminderUseCase: reminderUseCase,
		bot:             bot,
	}
}

func (b *botUseCase) HandleStartCommand(user *tgbotapi.User) (*keyboards.SelectionResult, error) {
	// Create or get user
	userEntity, err := b.userUseCase.CreateOrUpdateUser(
		user.ID,
		user.UserName,
		user.FirstName,
		user.LastName,
		"",
	)
	if err != nil {
		return nil, err
	}

	// Auto-detect language if not set
	if userEntity.Language == "" {
		if lang, supported := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); supported {
			err = b.userUseCase.UpdateUserLanguage(user.ID, lang)
			if err != nil {
				log.Printf("Failed to update user language: %v", err)
			} else {
				userEntity.Language = lang
			}
		}
	}

	text := ""
	markup := &tgbotapi.InlineKeyboardMarkup{}

	msg := tgbotapi.NewMessage(user.ID, "")
	if userEntity.Language == "" {
		// Ask for language if still not set
		text = "Select language / Оберіть мову"
		msg.ParseMode = "HTML"
		markup = keyboards.GetLanguageSelectionMarkup()
	} else {
		// Use language and show welcome
		text = keyboards.T(userEntity.Language).Welcome
		markup = keyboards.GetMainMenuMarkup(userEntity.Language)
	}

	return &keyboards.SelectionResult{Text: text, Markup: markup}, nil
}

func (b *botUseCase) HandleCallbackQuery(user *tgbotapi.User, callbackData string) (*keyboards.SelectionResult, error) {
	// Get user and selection
	userEntity, selection, err := b.getUserWithSelection(user.ID)
	if err != nil {
		return nil, err
	}

	// Handle language selection
	if keyboards.IsLanguageSelectionCallback(callbackData) {
		return b.handleLanguageSelection(user, callbackData, userEntity)
	}

	// Handle other callback types
	keyboardType := keyboards.GetKeyboardType(callbackData)
	switch keyboardType {
	case keyboards.Main:
		return b.handleMainMenu(user, userEntity)
	case keyboards.Reccurence:
		return b.handleRecurrenceSelection(user, callbackData, userEntity, selection)
	case keyboards.Time:
		return b.handleTimeSelection(user, callbackData, userEntity, selection)
	case keyboards.Week:
		return b.handleWeekSelection(user, callbackData, userEntity, selection)
	case keyboards.Message:
		return b.handleMessageSelection(user, callbackData, userEntity, selection)
	case keyboards.Reminders:
		// Check if this is a delete action
		if id, ok := keyboards.ParseDeleteReminderID(callbackData); ok {
			err := b.reminderUseCase.DeleteReminder(id, user.ID)
			if err != nil {
				log.Printf("Failed to delete reminder: %v", err)
			}
		}
		return b.handleRemindersList(user, userEntity)
	default:
		return nil, errors.NewDomainError("UNKNOWN_CALLBACK", "Unknown callback type", nil)
	}
}

func (b *botUseCase) HandleTextMessage(user *tgbotapi.User, text string) (*keyboards.SelectionResult, error) {
	// Get user and selection
	userEntity, selection, err := b.getUserWithSelection(user.ID)
	if err != nil {
		return nil, err
	}

	// Handle custom time input
	if selection.CustomTime && selection.SelectedTime == "" {
		return b.handleCustomTimeInput(user, text, userEntity, selection)
	}

	// Handle custom text input
	if selection.CustomText {
		return b.handleCustomTextInput(user, text, userEntity, selection)
	}

	return &keyboards.SelectionResult{Text: "I didn't understand that. Please use the menu buttons.", Markup: nil}, nil
}

func (b *botUseCase) ProcessKeyboardSelection(callbackQuery *tgbotapi.CallbackQuery) (*keyboards.SelectionResult, error) {
	log.Printf("'[%s] %s %s' selected '%s'",
		callbackQuery.From.UserName,
		callbackQuery.From.FirstName,
		callbackQuery.From.LastName,
		callbackQuery.Data)

	selectionResult, err := b.HandleCallbackQuery(callbackQuery.From, callbackQuery.Data)
	if err != nil {
		return nil, err
	}

	return selectionResult, nil
}

func (b *botUseCase) ProcessUserInput(message *tgbotapi.Message) (*keyboards.SelectionResult, error) {
	if message.IsCommand() {
		if message.Command() == "start" {
			return b.HandleStartCommand(message.From)
		}
	}

	if message.Text != "" {
		return b.HandleTextMessage(message.From, message.Text)
	}

	return &keyboards.SelectionResult{Text: "I didn't understand that. Please use the menu buttons.", Markup: nil}, nil
}

// Helper methods
func (b *botUseCase) getUserWithSelection(userID int64) (*entities.User, *entities.UserSelection, error) {
	user, err := b.userUseCase.GetUser(userID)
	if err != nil {
		return nil, nil, err
	}

	selection, err := b.userUseCase.GetUserSelection(userID)
	if err != nil {
		return nil, nil, err
	}

	return user, selection, nil
}

func (b *botUseCase) handleLanguageSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User) (*keyboards.SelectionResult, error) {
	lang := keyboards.ParseLanguageFromCallback(callbackData)
	err := b.userUseCase.UpdateUserLanguage(user.ID, lang)
	if err != nil {
		return nil, err
	}

	return &keyboards.SelectionResult{Text: keyboards.T(lang).Welcome, Markup: keyboards.GetMainMenuMarkup(lang)}, nil
}

func (b *botUseCase) handleMainMenu(user *tgbotapi.User, userEntity *entities.User) (*keyboards.SelectionResult, error) {
	return &keyboards.SelectionResult{Text: keyboards.T(userEntity.Language).Welcome, Markup: keyboards.GetMainMenuMarkup(userEntity.Language)}, nil
}

func (b *botUseCase) handleRecurrenceSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	result, err := keyboards.HandleRecurrenceTypeSelection(callbackData, userEntity, selection)
	if err != nil {
		return nil, err
	}

	err = b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}
	return result, nil
}

func (b *botUseCase) handleTimeSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	result := keyboards.HandleTimeSelection(callbackData, userEntity, selection)
	err := b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}
	return result, nil
}

func (b *botUseCase) handleWeekSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	result := keyboards.HandleWeekSelection(callbackData, &selection.WeekOptions, userEntity.Language)
	err := b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}
	return result, nil
}

func (b *botUseCase) handleMessageSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	result, completed := keyboards.HandleMessageSelection(callbackData, userEntity, selection)
	err := b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}
	if completed {
		_, err := b.reminderUseCase.CreateReminder(user.ID, selection)
		if err != nil {
			log.Printf("Failed to create reminder: %v", err)
		} else {
			result = keyboards.FormatReminderConfirmation(userEntity, selection)

			err = b.userUseCase.ClearUserSelection(user.ID)
			if err != nil {
				log.Printf("Failed to clear user selection: %v", err)
			}
		}
	}
	return result, nil
}

func (b *botUseCase) handleRemindersList(user *tgbotapi.User, userEntity *entities.User) (*keyboards.SelectionResult, error) {
	// Note: Reminder deletion is handled in the callback processing
	// This function just displays the reminders list

	// Get user reminders
	reminders, err := b.reminderUseCase.GetUserReminders(user.ID)
	if err != nil {
		return nil, err
	}

	return &keyboards.SelectionResult{Text: keyboards.FormatRemindersListText(reminders, userEntity.Language), Markup: keyboards.GetRemindersListMarkup(reminders, userEntity.Language)}, nil
}

func (b *botUseCase) handleCustomTimeInput(user *tgbotapi.User, text string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	selectionResult := keyboards.HandleCustomTimeSelection(text, &tgbotapi.MessageConfig{}, userEntity, selection)

	// Update user selection
	err := b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}

	return selectionResult, nil
}

func (b *botUseCase) handleCustomTextInput(user *tgbotapi.User, text string, userEntity *entities.User, selection *entities.UserSelection) (*keyboards.SelectionResult, error) {
	selectionResult, completed := keyboards.HandleCustomText(text, &tgbotapi.MessageConfig{}, userEntity, selection)

	// Update user selection
	err := b.userUseCase.UpdateUserSelection(user.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection: %v", err)
	}

	// If custom text was successful, create the reminder
	if completed {
		_, err := b.reminderUseCase.CreateReminder(user.ID, selection)
		if err != nil {
			log.Printf("Failed to create reminder: %v", err)
		} else {
			selectionResult = keyboards.FormatReminderConfirmation(userEntity, selection)

			// Clear user selection after successful reminder creation
			err = b.userUseCase.ClearUserSelection(user.ID)
			if err != nil {
				log.Printf("Failed to clear user selection: %v", err)
			}
		}
	}

	return selectionResult, nil
}
