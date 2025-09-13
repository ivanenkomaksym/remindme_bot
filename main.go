// bootstrap telegram bot
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/internal/keyboards"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/notifier"
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI

	welcomeMessage = "Welcome to the Reminder Bot!"
	reminderRepo   repositories.ReminderRepository
	userRepo       repositories.UserRepository
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in .env file")
	}
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true

	// Initialize repository based on STORAGE environment variable
	storageTypeStr := os.Getenv("STORAGE")
	if storageTypeStr == "" {
		storageTypeStr = "inmemory" // Default to in-memory storage
	}

	storageType, err := repositories.ToStorageType(storageTypeStr)
	if err != nil {
		log.Fatalf("Invalid storage type '%s': %v", storageTypeStr, err)
	}

	factory := repositories.NewReminderRepositoryFactory()
	reminderRepo = factory.CreateRepository(storageType)
	log.Printf("Initialized %s storage repository", storageType.String())

	// Initialize user repository
	userFactory := repositories.NewUserRepositoryFactory()
	userRepo = userFactory.CreateRepository(storageType)
	log.Printf("Initialized %s user repository", storageType.String())

	// Start the reminder notification loop
	go notifier.StartReminderNotifier(reminderRepo, bot)

	// Get the WEBHOOK_URL from environment variables
	// This will be the URL of your deployed Cloud Run service + the webhook path
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable not set. This should be your Cloud Run service URL including the path (e.g., https://<service-url>/telegram-webhook).")
	}

	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Failed to create webhook config: %v", err)
	}
	// Use bot.Request to send the WebhookConfig to Telegram
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	// Get webhook info to confirm it's set and check for any errors from Telegram's side
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatalf("Failed to get webhook info: %v", err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram last webhook error: %s", info.LastErrorMessage)
	}
	log.Printf("Webhook set to: %s (pending: %t)", info.URL, info.PendingUpdateCount > 0)

	// Get the port from environment variables, default to 8080 for Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

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
		handleUpdate(update)

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

func handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		processKeyboardSelection(update.CallbackQuery)
		return
	}

	if update.Message != nil {
		processUserInput(update.Message)
		return
	}
}

func processKeyboardSelection(callbackQuery *tgbotapi.CallbackQuery) bool {
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
	userModel, userState := userRepo.CreateOrUpdateUserWithState(user.ID,
		user.UserName,
		user.FirstName,
		user.LastName,
		"")

	if userModel.Language == "" {
		if lang, supported := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); supported {
			userModel.Language = lang
			userRepo.UpdateUserLanguage(userModel.Id, lang)
		}
	}

	// Language selection
	if keyboards.IsLanguageSelectionCallback(text) {
		lang := keyboards.ParseLanguageFromCallback(text)
		userModel.Language = lang
		userRepo.UpdateUserLanguage(userModel.Id, lang)
		msg.Text = keyboards.T(lang).Welcome
		markup = keyboards.GetMainMenuMarkup(userModel.Language)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = markup
		bot.Send(msg)
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
		userRepo.UpdateUserState(userModel.Id, userState)
	case keyboards.Time:
		markup = keyboards.HandleTimeSelection(text, &msg, userModel, userState)
		userRepo.UpdateUserState(userModel.Id, userState)
	case keyboards.Week:
		markup = keyboards.HandleWeekSelection(text, &msg, &userState.WeekOptions, userModel.Language)
		userRepo.UpdateUserState(userModel.Id, userState)
	case keyboards.Message:
		messageMarkup, completed := keyboards.HandleMessageSelection(text, &msg, userModel, userState)
		userRepo.UpdateUserState(userModel.Id, userState)

		// If message selection was successful, create the reminder
		if completed {
			if !handleReminderCreation(userModel, userState, &msg) {
				return false
			}
		}

		markup = messageMarkup
	case keyboards.Reminders:
		// Show or update the reminders list, and handle deletions
		if id, ok := keyboards.ParseDeleteReminderID(text); ok {
			_ = reminderRepo.DeleteReminder(id, userModel.Id)
		}
		userRems := reminderRepo.GetRemindersByUser(userModel.Id)
		msg.Text = keyboards.FormatRemindersListText(userRems, userModel.Language)
		markup = keyboards.GetRemindersListMarkup(userRems, userModel.Language)
	}
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	msg.ParseMode = "HTML"
	bot.Send(msg)
	return true
}

func processUserInput(message *tgbotapi.Message) bool {
	user := message.From
	text := message.Text

	if message.IsCommand() {
		log.Printf("'[%s] %s %s' started chat", user.UserName, user.FirstName, user.LastName)

		if message.Command() == "start" {
			userModel, _ := userRepo.CreateOrUpdateUserWithState(user.ID, user.UserName, user.FirstName, user.LastName, "")
			// Try to auto-detect language if not set yet
			if userModel.Language == "" {
				if lang, ok := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); ok {
					userModel.Language = lang
					userRepo.UpdateUserLanguage(userModel.Id, lang)
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
			bot.Send(msg)
		}
	} else if text != "" {
		// Handle custom time input or custom message input
		userModel, userState := userRepo.CreateOrUpdateUserWithState(user.ID, user.UserName, user.FirstName, user.LastName, "")

		// Handle language detection externally
		if userModel.Language == "" {
			if lang, supported := keyboards.MapTelegramLanguageCodeToSupported(user.LanguageCode); supported {
				userModel.Language = lang
				userRepo.UpdateUserLanguage(userModel.Id, lang)
			}
		}

		msg := tgbotapi.NewMessage(
			message.Chat.ID,
			"", // Text will be set later
		)

		if userState.CustomTime && userState.SelectedTime == "" {
			msg.ReplyMarkup = keyboards.HadleCustomTimeSelection(text, &msg, userModel, userState)
			userRepo.UpdateUserState(userModel.Id, userState)
		} else if userState.CustomText {
			markup, completed := keyboards.HadleCustomText(text, &msg, userModel, userState)
			msg.ReplyMarkup = markup
			userRepo.UpdateUserState(userModel.Id, userState)

			// If custom text was successful, create the reminder
			if completed {
				if !handleReminderCreation(userModel, userState, &msg) {
					return false
				}
			}
		}

		msg.ParseMode = "HTML"
		bot.Send(msg)
	}

	return true
}

func handleReminderCreation(user *models.User, userState *types.UserSelectionState, msg any) bool {
	err := createReminder(user, userState)
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

	userRepo.ClearUserState(user.Id)

	return true
}

func createReminder(user *models.User, userState *types.UserSelectionState) error {
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
