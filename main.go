// bootstrap telegram bot
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ivanenkomaksym/remindme_bot/internal/keyboards"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI

	welcomeMessage = "Welcome to the Reminder Bot!"
)

func buildMainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Daily.String(), models.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Weekly.String(), models.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Monthly.String(), models.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Interval.String(), models.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Custom.String(), models.Custom.String()),
		),
	)
}

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
		user := update.CallbackQuery.From
		text := update.CallbackQuery.Data

		log.Printf("'[%s] %s %s' selected '%s'", user.UserName, user.FirstName, user.LastName, text)

		msg := tgbotapi.NewEditMessageText(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			"", // Текст буде встановлено нижче
		)
		var markup *tgbotapi.InlineKeyboardMarkup

		// Check if this is a time selection callback
		keyboardType := keyboards.GetKeyboardType(update.CallbackQuery.Data)
		switch keyboardType {
		case keyboards.Reccurence:
			recurrenceType, err := models.ToRecurrenceType(update.CallbackQuery.Data)
			if err != nil {
				log.Printf("Failed to resolve selected recurrence type: %v", err)
				return
			}
			markup = handleRecurrenceTypeSelection(recurrenceType, &msg)
		case keyboards.Time:
			markup = handleTimeSelection(update, &msg)
		}
		if markup != nil {
			msg.ReplyMarkup = markup
		}

		msg.ParseMode = "HTML"
		bot.Send(msg)
	}

	if update.Message != nil {
		user := update.Message.From

		if update.Message.IsCommand() {
			log.Printf("'[%s] %s %s' started chat", user.UserName, user.FirstName, user.LastName)

			if update.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = buildMainMenu()
				bot.Send(msg)
			}
		} else if update.Message.Text != "" {
			// Handle custom time input (e.g., "14:30")
			if isValidTimeFormat(update.Message.Text) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Custom time '%s' accepted! You will receive daily reminders at this time.", update.Message.Text))
				msg.ParseMode = "HTML"
				// Clear markup for confirmation
				bot.Send(msg)
			}
		}
	}
}

func handleRecurrenceTypeSelection(recurrenceType models.RecurrenceType, msg *tgbotapi.EditMessageTextConfig) *tgbotapi.InlineKeyboardMarkup {
	switch recurrenceType {
	case models.Daily:
		msg.Text = "Select time for daily reminders:"
		menu := keyboards.GetHourRangeMarkup()
		return &menu
	case models.Weekly:
		msg.Text = "Select time for weekly reminders:"
		menu := keyboards.GetHourRangeMarkup()
		return &menu
	case models.Monthly:
		msg.Text = "Select time for monthly reminders:"
		menu := keyboards.GetHourRangeMarkup()
		return &menu
	case models.Interval:
		msg.Text = "Select time for interval reminders:"
		menu := keyboards.GetHourRangeMarkup()
		return &menu
	case models.Custom:
		msg.Text = "Please type your custom time in HH:MM format (e.g., 14:05):"
		menu := keyboards.GetHourRangeMarkup()
		return &menu
	}

	return nil
}

func handleTimeSelection(update tgbotapi.Update, msg *tgbotapi.EditMessageTextConfig) *tgbotapi.InlineKeyboardMarkup {
	callbackData := update.CallbackQuery.Data

	switch {
	case callbackData == "back_to_main":
		// User wants to go back to main menu
		msg.Text = "Select reminder frequency:"
		menu := buildMainMenu()
		return &menu

	case callbackData == "back_to_hour_range":
		// User wants to go back to hour range selection
		msg.Text = "Select time for your reminders:"
		menu := keyboards.GetHourRangeMarkup()
		return &menu

	case strings.Contains(callbackData, keyboards.CallbackPrefixHourRange):
		// User selected a 4-hour range, show 1-hour ranges
		startHour := 0
		fmt.Sscanf(callbackData[len(keyboards.CallbackPrefixHourRange):], "%d", &startHour)
		msg.Text = fmt.Sprintf("Select hour within %02d:00-%02d:00:", startHour, (startHour+4)%24)
		menu := keyboards.GetMinuteRangeMarkup(startHour)
		return &menu

	case strings.Contains(callbackData, keyboards.CallbackPrefixMinuteRange):
		// User selected a 1-hour range, show 15-minute intervals
		startHour := 0
		fmt.Sscanf(callbackData[len(keyboards.CallbackPrefixMinuteRange):], "%d", &startHour)
		msg.Text = fmt.Sprintf("Select time within %02d:00-%02d:00:", startHour, (startHour+1)%24)
		menu := keyboards.GetSpecificTimeMarkup(startHour)
		return &menu

	case strings.Contains(callbackData, keyboards.CallbackPrefixSpecificTime):
		// User selected a specific time, show confirmation
		timeStr := callbackData[len(keyboards.CallbackPrefixSpecificTime):]
		msg.Text = fmt.Sprintf("Reminder set for %s! You will receive daily reminders at this time.", timeStr)
		return nil

	case strings.Contains(callbackData, keyboards.CallbackPrefixCustom):
		// User wants custom time input
		msg.Text = "Please type your custom time in HH:MM format (e.g., 14:30):"
		return nil
	}

	return nil
}

// isValidTimeFormat checks if the input string is a valid time format (HH:MM)
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return false
	}

	hour := timeStr[:2]
	minute := timeStr[3:]

	// Check if hour and minute are valid numbers
	if hour < "00" || hour > "23" || minute < "00" || minute > "59" {
		return false
	}

	return true
}
