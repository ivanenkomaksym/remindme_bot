package keyboards

import (
	"log"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleNlpTextInputCallback handles the NLP text input callback and sets up the user selection
func HandleNlpTextInputCallback(userEntity *entities.User, updateUserSelection func(int64, *entities.UserSelection) error) (*SelectionResult, error) {
	s := T(userEntity.Language)

	// Create a new selection and set it to expect custom text input
	// We'll use the CustomText flag to identify that user should input text
	// and use ReminderMessage as a marker for NLP mode
	selection := entities.NewUserSelection()
	selection.SetCustomText()
	selection.ReminderMessage = "NLP_MODE" // Use this as a marker for NLP processing

	err := updateUserSelection(userEntity.ID, selection)
	if err != nil {
		log.Printf("Failed to update user selection for NLP: %v", err)
		return &SelectionResult{
			Text:   s.MsgParsingFailed,
			Markup: GetNavigationMenuMarkup(userEntity.Language),
		}, nil
	}

	text := s.NlpMenuTitle + "\n\n" + s.NlpInstructions + "\n\n" + s.NlpExamples + "\n\n" + s.NlpEnterText

	// Create a back button markup
	backButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackSetup),
		),
	)

	return &SelectionResult{
		Text:   text,
		Markup: &backButton,
	}, nil
}

// NLPService interface for parsing reminder text
type NLPService interface {
	ParseReminderText(text string, userTimezone string, userLanguage string) (*entities.UserSelection, error)
}

// HandleNlpTextProcessing processes the NLP text input and creates a reminder
func HandleNlpTextProcessing(
	text string,
	userEntity *entities.User,
	nlpService NLPService,
	createReminder func(int64, *entities.UserSelection) (*entities.Reminder, error),
	clearUserSelection func(int64) error,
) (*SelectionResult, error) {
	s := T(userEntity.Language)

	// Get user's timezone for NLP processing
	timezone := "UTC"
	if userEntity.GetLocation() != nil {
		timezone = userEntity.GetLocation().String()
	}

	// Use NLP service to parse the text
	nlpSelection, err := nlpService.ParseReminderText(text, timezone, userEntity.Language)
	if err != nil {
		log.Printf("Failed to parse NLP text: %v", err)
		return &SelectionResult{
			Text: "❌ " + s.MsgParsingFailed + "\n\n" + s.NlpEnterText,
			Markup: &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackSetup)},
				},
			},
		}, nil
	}

	// Create the reminder using the parsed selection
	_, err = createReminder(userEntity.ID, nlpSelection)
	if err != nil {
		log.Printf("Failed to create NLP reminder: %v", err)
		return &SelectionResult{
			Text:   "❌ Error creating reminder. Please try again.",
			Markup: GetNavigationMenuMarkup(userEntity.Language),
		}, nil
	}

	// Clear user selection after successful reminder creation
	err = clearUserSelection(userEntity.ID)
	if err != nil {
		log.Printf("Failed to clear user selection: %v", err)
	}

	// Format confirmation message
	confirmationResult := FormatReminderConfirmation(userEntity, nlpSelection)
	return confirmationResult, nil
}
