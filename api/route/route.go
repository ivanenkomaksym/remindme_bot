package route

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/bootstrap"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"
	"github.com/ivanenkomaksym/remindme_bot/models"
	"github.com/ivanenkomaksym/remindme_bot/repositories"
	"github.com/ivanenkomaksym/remindme_bot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Setup(app *bootstrap.Application) {
	addr := app.Env.ServerAddress

	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Define the webhook endpoint that Telegram will send updates to
	// It's good practice to use a non-root path for webhooks.
	mux.HandleFunc("/telegram-webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var update tgbotapi.Update
		// Decode the JSON request body into a tgbotapi.Update struct
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Printf("ERROR: Could not decode update: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Process the update in a non-blocking way if possible,
		// but for simple bots, direct handling is fine.
		handleUpdate(app, update)

		// Respond with 200 OK to Telegram immediately
		// This acknowledges receipt of the update and prevents Telegram from retrying.
		w.WriteHeader(http.StatusOK)
	})

	// Add a health check endpoint for Cloud Run
	// Cloud Run sends requests to the root path by default for health checks.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting HTTP server on %s", addr)
	// Start the HTTP server. This will block indefinitely, serving requests.
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func handleUpdate(app *bootstrap.Application, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		processKeyboardSelection(app, update.CallbackQuery)
		return
	}

	if update.Message != nil {
		processUserInput(app, update.Message)
		return
	}
}

func processKeyboardSelection(app *bootstrap.Application, callbackQuery *tgbotapi.CallbackQuery) bool {
	user := callbackQuery.From
	text := callbackQuery.Data

	log.Printf("'[%s] %s %s' selected '%s'", user.UserName, user.FirstName, user.LastName, text)

	msg := tgbotapi.NewEditMessageText(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		"", // Text will be set later
	)
	var markup *tgbotapi.InlineKeyboardMarkup

	// load or init per-user selection state
	userModel, userState := app.UserRepo.CreateOrUpdateUserWithState(user.ID,
		user.UserName,
		user.FirstName,
		user.LastName,
		"")

	if userModel.Language == "" {
		if lang, supported := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); supported {
			userModel.Language = lang
			app.UserRepo.UpdateUserLanguage(userModel.Id, lang)
		}
	}

	// Language selection
	if keyboards.IsLanguageSelectionCallback(text) {
		lang := keyboards.ParseLanguageFromCallback(text)
		userModel.Language = lang
		app.UserRepo.UpdateUserLanguage(userModel.Id, lang)
		msg.Text = keyboards.T(lang).Welcome
		markup = keyboards.GetMainMenuMarkup(userModel.Language)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = markup
		app.Bot.Send(msg)
		return true
	}

	// Check if this is a time selection callback
	keyboardType := keyboards.GetKeyboardType(text)
	switch keyboardType {
	case keyboards.Main:
		msg.Text = keyboards.T(userModel.Language).Welcome
		markup = keyboards.GetMainMenuMarkup(userModel.Language)
	case keyboards.Reccurence:
		m, err := keyboards.HandleRecurrenceTypeSelection(text, &msg, userModel, userState)
		if err != nil {
			log.Printf("Failed to resolve selected recurrence type: %v", err)
			return false
		}
		markup = m
	case keyboards.Time:
		markup = keyboards.HandleTimeSelection(text, &msg, userModel, userState)
	case keyboards.Week:
		markup = keyboards.HandleWeekSelection(text, &msg, &userState.WeekOptions, userModel.Language)
	case keyboards.Message:
		messageMarkup, completed := keyboards.HandleMessageSelection(text, &msg, userModel, userState)

		// If message selection was successful, create the reminder
		if completed {
			if !handleReminderCreation(app.ReminderRepo, userModel, userState, &msg) {
				return false
			}
		}

		markup = messageMarkup
	case keyboards.Reminders:
		// Show or update the reminders list, and handle deletions
		if id, ok := keyboards.ParseDeleteReminderID(text); ok {
			_ = app.ReminderRepo.DeleteReminder(id, userModel.Id)
		}
		userRems := app.ReminderRepo.GetRemindersByUser(userModel.Id)
		msg.Text = keyboards.FormatRemindersListText(userRems, userModel.Language)
		markup = keyboards.GetRemindersListMarkup(userRems, userModel.Language)
	}
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	app.UserRepo.UpdateUserState(userModel.Id, userState)

	msg.ParseMode = "HTML"
	app.Bot.Send(msg)
	return true
}

func processUserInput(app *bootstrap.Application, message *tgbotapi.Message) bool {
	user := message.From
	text := message.Text

	if message.IsCommand() {
		log.Printf("'[%s] %s %s' started chat", user.UserName, user.FirstName, user.LastName)

		if message.Command() == "start" {
			userModel, _ := app.UserRepo.CreateOrUpdateUserWithState(user.ID, user.UserName, user.FirstName, user.LastName, "")
			// Try to auto-detect language if not set yet
			if userModel.Language == "" {
				if lang, ok := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); ok {
					userModel.Language = lang
					app.UserRepo.UpdateUserLanguage(userModel.Id, lang)
				}
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "")
			if userModel.Language == "" {
				// Ask for language if still not set
				msg.Text = "Select language / Оберіть мову"
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = keyboards.GetLanguageSelectionMarkup()
			} else {
				// Use language and show welcome
				msg.Text = keyboards.T(userModel.Language).Welcome
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = keyboards.GetMainMenuMarkup(userModel.Language)
			}
			app.Bot.Send(msg)
		}
	} else if text != "" {
		// Handle custom time input or custom message input
		userModel, userState := app.UserRepo.CreateOrUpdateUserWithState(user.ID, user.UserName, user.FirstName, user.LastName, "")

		// Handle language detection externally
		if userModel.Language == "" {
			if lang, supported := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); supported {
				userModel.Language = lang
				app.UserRepo.UpdateUserLanguage(userModel.Id, lang)
			}
		}

		msg := tgbotapi.NewMessage(
			message.Chat.ID,
			"", // Text will be set later
		)

		if userState.CustomTime && userState.SelectedTime == "" {
			msg.ReplyMarkup = keyboards.HadleCustomTimeSelection(text, &msg, userModel, userState)
			app.UserRepo.UpdateUserState(userModel.Id, userState)
		} else if userState.CustomText {
			markup, completed := keyboards.HadleCustomText(text, &msg, userModel, userState)
			msg.ReplyMarkup = markup
			app.UserRepo.UpdateUserState(userModel.Id, userState)

			// If custom text was successful, create the reminder
			if completed {
				if !handleReminderCreation(app.ReminderRepo, userModel, userState, &msg) {
					return false
				}
				app.UserRepo.ClearUserState(user.ID)
			}
		}

		msg.ParseMode = "HTML"
		app.Bot.Send(msg)
	}

	return true
}

func handleReminderCreation(reminderRepo repositories.ReminderRepository, user *models.User, userState *types.UserSelectionState, msg any) bool {
	err := createReminder(reminderRepo, user, userState)
	if err != nil {
		log.Printf("Could not create reminder due to error: %s", err)
		return false
	}

	// Set the confirmation message text based on the message type
	text, keyboard := keyboards.FormatReminderConfirmation(user, userState)
	switch m := msg.(type) {
	case *tgbotapi.EditMessageTextConfig:
		m.Text = text
		m.ReplyMarkup = keyboard
	case *tgbotapi.MessageConfig:
		m.Text = text
		m.ReplyMarkup = keyboard
	}

	return true
}

func createReminder(reminderRepo repositories.ReminderRepository, user *models.User, userState *types.UserSelectionState) error {
	switch userState.RecurrenceType {
	case models.Daily:
		reminderRepo.CreateDailyReminder(userState.SelectedTime, *user, userState.ReminderMessage)
	case models.Weekly:
		// Convert week options to time.Weekday slice
		var daysOfWeek []time.Weekday
		for i, selected := range userState.WeekOptions {
			if selected {
				daysOfWeek = append(daysOfWeek, time.Weekday(i))
			}
		}
		reminderRepo.CreateWeeklyReminder(daysOfWeek, userState.SelectedTime, *user, userState.ReminderMessage)
	case models.Monthly:
		// For monthly, we'll use the 1st of each month for now
		// This could be enhanced to allow user to select specific days
		daysOfMonth := []int{1}
		reminderRepo.CreateMonthlyReminder(daysOfMonth, userState.SelectedTime, *user, userState.ReminderMessage)
	case models.Interval:
		// For interval, treat as daily for now
		reminderRepo.CreateDailyReminder(userState.SelectedTime, *user, userState.ReminderMessage)
	case models.Custom:
		// For custom, treat as daily for now
		reminderRepo.CreateDailyReminder(userState.SelectedTime, *user, userState.ReminderMessage)
	}

	return nil
}
