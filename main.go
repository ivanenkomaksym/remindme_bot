// bootstrap telegram bot
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/internal/keyboards"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI

	welcomeMessage = "Welcome to the Reminder Bot!"

	// Global repository instance
	reminderRepo repositories.ReminderRepository
)

// in-memory, per-user selection state
var (
	userSelectionsMu     sync.RWMutex
	userSelectionsByUser = map[int64]*types.UserSelectionState{}
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
	userState := getUserStateWithUser(callbackQuery.From)

	// Check if this is a time selection callback
	keyboardType := keyboards.GetKeyboardType(callbackQuery.Data)
	switch keyboardType {
	case keyboards.Main:
		msg.Text = welcomeMessage
		markup = keyboards.GetMainMenuMarkup()
	case keyboards.Reccurence:
		m, err := keyboards.HandleRecurrenceTypeSelection(callbackQuery.Data, &msg, userState)
		if err != nil {
			log.Printf("Failed to resolve selected recurrence type: %v", err)
			return false
		}
		markup = m
	case keyboards.Time:
		markup = keyboards.HandleTimeSelection(callbackQuery.Data, &msg, userState)
	case keyboards.Week:
		userSelectionsMu.Lock()
		markup = keyboards.HandleWeekSelection(callbackQuery.Data, &msg, &userState.WeekOptions)
		userSelectionsMu.Unlock()
	case keyboards.Message:
		userSelectionsMu.Lock()
		markup = keyboards.HandleMessageSelection(callbackQuery.Data, &msg, userState, reminderRepo)
		userSelectionsMu.Unlock()
	}
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	msg.ParseMode = "HTML"
	bot.Send(msg)
	return true
}

func processUserInput(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if message.IsCommand() {
		log.Printf("'[%s] %s %s' started chat", user.UserName, user.FirstName, user.LastName)

		if message.Command() == "start" {
			msg := tgbotapi.NewMessage(message.Chat.ID, welcomeMessage)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = keyboards.GetMainMenuMarkup()
			bot.Send(msg)
		}
	} else if text != "" {
		// Handle custom time input or custom message input
		userState := getUserStateWithUser(message.From)

		msg := tgbotapi.NewMessage(
			message.Chat.ID,
			"", // Text will be set later
		)

		if userState.CustomTime && userState.SelectedTime == "" {
			msg.ReplyMarkup = keyboards.HadleCustomTimeSelection(text, &msg, userState)
		} else if userState.CustomText {
			msg.ReplyMarkup = keyboards.HadleCustomText(text, &msg, userState, reminderRepo)
		}

		msg.ParseMode = "HTML"
		bot.Send(msg)
	}
}

func getUserStateWithUser(user *tgbotapi.User) *types.UserSelectionState {
	userSelectionsMu.Lock()
	userState, ok := userSelectionsByUser[user.ID]
	if !ok {
		userState = &types.UserSelectionState{
			User: models.User{
				Id:        user.ID,
				UserName:  user.UserName,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
			WeekOptions: [7]bool{false, false, false, false, false, false, false},
		}
		userSelectionsByUser[user.ID] = userState
	} else {
		// Update user info in case it changed
		userState.User = models.User{
			Id:        user.ID,
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}
	userSelectionsMu.Unlock()
	return userState
}
